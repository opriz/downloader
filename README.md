# golang downloader

now support http/https fast download

## notice

some sites may require for headers to download resources, not support now 

## install 

```
    git clone https://github.com/opriz/downloader.git   
    cd downloader    
    go build  
```

## usage
```
./downloader -h
    -n, --file-name string   new file name (default the last part of url)   
    -p, --parralel int       download parralels (default 5)   
    -u, --url string         download url   
```
