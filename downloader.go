package main

import(
    dwler "downloader"
)

type Parser interface {
    Start() chan string
    AfterFinished(url string)
}

type Downloader interface{
    Download(url string) []byte
}

type ImgResource interface{
    GetUrl() string
}

type Content struct{
    Content []byte
    Resource ImgResource
}

var DownloadWorker chan *Downloader

func MakeDownloaderWorkers() {
    var worker_factory func() Downloader
    switch Config.UseProxy{
    case true:
        worker_factory = ProxyDownloaderFactory
        break
    case false:
        worker_factory = DirectDownloaderFactory
        break
    }
    DownloadWorker = make(chan *Downloader, CONCURENT_DOWNLOADS)
    for i:=0;i<CONCURENT_DOWNLOADS;i++{
        downloader := worker_factory()
        DownloadWorker <- &downloader
    }
}

func ProxyDownloaderFactory() Downloader{
    downloader := dwler.MakePDownloader(Config.Proxy, Config.Timeout)
    return downloader
}

func DirectDownloaderFactory() Downloader{
    downloader := dwler.MakeDirectDownloader(Config.Timeout)
    return downloader
}

func Download_raw(img_resource ImgResource, td TumblrDownloader){
    worker := *(<-DownloadWorker)
    defer func(){
        DownloadWorker <- &worker
    }()
    content := worker.Download(img_resource.GetUrl())
    if len(content) > 0{
        result := Content{
            Content:content,
            Resource:img_resource,
        }
        td.ContenChan <- result
    }else{
        td.UrlChan <- img_resource
    }
}
