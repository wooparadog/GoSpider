package main

import(
)

var exit_signal chan int

func main(){
    ParseConfig()
    recorder := CreateRecorders()
    MakeDownloaderWorkers()
    for _, tumblr_source := range Config.TumblrSources{
        td := MakeTumblrDownloader(tumblr_source.Name, tumblr_source.Suffix, tumblr_source.Url, recorder)
        go td.Start()
    }
    a:=make(chan int)
    <-a
}
