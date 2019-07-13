package main

import (
	"net/http"
	"time"

	"download/downloader"
	"download/logger"

	flag "github.com/spf13/pflag"
)

type Config struct {
	URL      string
	FileName string
	LogLevel string
	Parralel int64
}

func main() {
	var config Config
	flag.StringVarP(&config.URL, "url", "u", "", "download url")
	flag.StringVarP(&config.FileName, "file-name", "n", "", "new file name")
	flag.Int64VarP(&config.Parralel, "parralel", "p", 5, "download parralels")
	flag.StringVarP(&config.LogLevel, "log-level", "l", "Debug", "log level<Debug,Info,Warn,Error,Panic,Fatal>")
	flag.Parse()

	logger := logger.NewZapLogger(config.LogLevel)

	if config.URL == "" {
		logger.Errorf("no URL")
		return
	}

	downloader := downloader.NewHTTPdownloader(
		logger,
		&http.Client{},
		config.Parralel,
	)

	startTime := time.Now()
	err := downloader.Download(config.URL, config.FileName)
	if err != nil {
		logger.Errorf("Download failed, %v", err)
		return
	}
	endTime := time.Now()
	cost := endTime.Sub(startTime)
	logger.Infof("Download done, time cost: %v", cost)
}
