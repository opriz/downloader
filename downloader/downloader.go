package downloader

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"download/logger"
)

type Downloader interface {
	Download(url string) error
}

type downloadBlock struct {
	start int64
	stop  int64
}

type HTTPdownloader struct {
	logger   logger.Logger
	client   *http.Client
	parralel int64
}

func NewHTTPdownloader(logger logger.Logger, httpClient *http.Client, parralel int64) *HTTPdownloader {
	return &HTTPdownloader{
		logger:   logger,
		client:   httpClient,
		parralel: parralel,
	}
}

func (d *HTTPdownloader) Download(rawurl string, filename string) error {
	blocks, err := d.makeBlocks(rawurl)
	if err != nil {
		d.logger.Errorf("HTTPdownloader: Download: failed, rawurl=%v, filename=%v, %v", rawurl, filename, err)
		return err
	}

	wg := &sync.WaitGroup{}
	errChan := make(chan error)
	infoChan := make(chan *downloadBlock)
	filename = d.resolveFilename(rawurl, filename)
	d.logger.Debugf("HTTPdownloader: Download: filename=%v", filename)

	file, err := os.Create(filename)
	if err != nil {
		d.logger.Errorf("HTTPdownloader: Download: file Create failed, %v", err)
		return err
	}
	defer file.Close()

	for _, block := range blocks {
		wg.Add(1)
		d.logger.Debugf("HTTPdownloader: Download: add block, %+v", *block)
		go d.downloadBlock(rawurl, block, file, infoChan, errChan, wg)
	}

	for i := int64(0); i < d.parralel; i++ {
		select {
		case err := <-errChan:
			d.logger.Errorf("HTTPdownloader: Download: get error: %v", err)
		case block := <-infoChan:
			d.logger.Infof("HTTPdownloader: Download: finish: %v-%v", block.start, block.stop)
		}
	}
	wg.Wait()
	return nil
}

func (d *HTTPdownloader) makeBlocks(url string) ([]*downloadBlock, error) {
	resp, err := d.client.Head(url)
	if err != nil {
		// resp status 不 ok 时err？
		d.logger.Errorf("HTTPdownloader: makeBlocks: Head error, %v", err)
		return nil, err
	}
	d.logger.Debugf("HTTPdownloader: makeBlocks: ContentLength=%v", resp.ContentLength)
	normalBlockLen := resp.ContentLength / d.parralel
	lastBlockLen := resp.ContentLength - normalBlockLen*d.parralel + normalBlockLen

	//blocks := make([]*downloadBlock, d.parralel)
	blocks := []*downloadBlock{}
	begin := int64(0)
	for i := int64(0); i < d.parralel-1; i++ {
		b := &downloadBlock{start: begin, stop: begin + normalBlockLen - 1}
		blocks = append(blocks, b)
		d.logger.Debugf("HTTPdownloader: makeBlocks: block: %v-%v, length=%v", b.start, b.stop, b.stop-b.start+1)
		begin = begin + normalBlockLen
	}
	lastBlock := &downloadBlock{start: begin, stop: begin + lastBlockLen - 1}
	d.logger.Debugf("HTTPdownloader: makeBlocks: block: %v-%v", lastBlock.start, lastBlock.stop)
	blocks = append(blocks, lastBlock)
	return blocks, nil
}

func (d *HTTPdownloader) downloadBlock(
	url string,
	block *downloadBlock,
	file *os.File,
	infoChan chan *downloadBlock,
	errChan chan error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		d.logger.Errorf("HTTPdownloader: downloadBlock: http NewRequest failed, %v", err)
		errChan <- err
		return
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", block.start, block.stop))
	resp, err := d.client.Do(req)
	if err != nil {
		d.logger.Errorf("HTTPdownloader: downloadBlock: block=%+v, failed, %v", block, err)
		errChan <- err
		return
	}
	d.logger.Debugf("HTTPdownloader: downloadBlock: response status=%v, downloading...", resp.Status)

	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(resp.Body)
	if err != nil {
		d.logger.Errorf("HTTPdownloader: downloadBlock: buf ReadFrom failed, %v", err)
		errChan <- err
		return
	}
	d.logger.Debugf("HTTPdownloader: downloadBlock: buf ReadFrom: n=%v", n)

	_, err = file.WriteAt(buf.Bytes(), block.start)
	if err != nil {
		d.logger.Errorf("HTTPdownloader: downloadBlock: file WriteAt failed, %v", err)
		errChan <- err
		return
	}
	infoChan <- block
}

func (d *HTTPdownloader) resolveFilename(rawurl, filename string) string {
	if filename != "" {
		return filename
	}

	words := strings.Split(rawurl, "/")
	return words[len(words)-1]
}
