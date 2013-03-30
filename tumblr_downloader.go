package main

import(
    "os"
    "log"
    "io/ioutil"
    "path"
    "strings"
    "regexp"
    "time"
    "errors"
)


const (
    FullResolution = iota
    OriginalResolution = iota
)

type TumblrImg struct{
    id string
    url string
    state int
    rss_img_size string
}

func (self *TumblrImg) GetUrl() string{
    url := ""
    switch(self.state){
    case OriginalResolution:
        url = self.url
        break
    default:
        url = strings.Replace(self.url, self.rss_img_size, "_1280", -1)
    }
    return url
}

type TumblrDownloader struct{
    recorder Recorder
    UrlChan chan ImgResource
    ContenChan chan Content
    img_id_finder *regexp.Regexp
    name string
    rss_img_size string
    rss_addr string
}

func MakeTumblrDownloader(name, rss_img_size, rss_addr string, recorder Recorder) TumblrDownloader{
    img_id_finder,_ := regexp.Compile(`tumblr_(.*?)_`)
    downloader := TumblrDownloader{
        recorder:recorder,
        name:name,
        rss_img_size:rss_img_size,
        rss_addr:rss_addr,
        ContenChan:make(chan Content),
        UrlChan:make(chan ImgResource),
        img_id_finder: img_id_finder,
    }
    return downloader
}

func (self *TumblrDownloader) get_img_id(url string) (string, error){
    filename := path.Base(url)
    matches := self.img_id_finder.FindStringSubmatch(filename)
    if len(matches) > 0{
        return matches[1], nil
    }
    log.Printf("WARNING: Fail to Match id for %s", url)
    return "", errors.New("Fail to Match Id: " + url)
}

func (self *TumblrDownloader) AfterFinished(){
    dir, _ := os.Getwd()
    dir = path.Join(dir, "Downloads", self.name)
    _, err := os.Stat(dir)
    _, ok := err.(*os.PathError)
    if ok{
        log.Printf("make dir : %s", dir)
        os.MkdirAll(dir, 0775)
    }
    for {
        content := <- self.ContenChan
        id, err := self.get_img_id(content.Resource.GetUrl())
        if err == nil{
            filename := path.Base(content.Resource.GetUrl())
            filename = path.Join(dir, filename)
            if strings.Contains(string(content.Content), "AccessDenied"){
                log.Printf("Access Denied Error, Try Use Original Link : %s", content.Resource.(*TumblrImg).url)
                self.AddUrl(&TumblrImg{
                    id:id,
                    url:content.Resource.(*TumblrImg).url,
                    state:OriginalResolution,
                    rss_img_size:self.rss_img_size,
                })
            }else{
                ioutil.WriteFile(filename, content.Content, 0600)
                self.recorder.MarkAsFinished(id)
                log.Printf("New Download Saved : %s", filename)
            }
        }
    }
}

func (self *TumblrDownloader) AddUrl(tumblr_img ImgResource){
    self.UrlChan <- tumblr_img
}

func (self *TumblrDownloader) GetWaitingUrl() (bool, string){
    tumblr_img, ok := <- self.UrlChan
    return ok, tumblr_img.GetUrl()
}

func (self *TumblrDownloader) ProcessUrl(url string){
    id, err := self.get_img_id(url)
    if err == nil {
        if !self.recorder.HasFinished(id){
            tumblr_img := TumblrImg{
                id:id,
                url:url,
                state:FullResolution,
                rss_img_size:self.rss_img_size,
            }
            self.AddUrl(&tumblr_img)
        }else{
            //log.Printf("Already Downloaded : %s", id)
        }
    }
}

func (self *TumblrDownloader) download_rss() string{
    p_downloader := *(<- DownloadWorker)
    defer func(){
        DownloadWorker <- &p_downloader
    }()
    return string(p_downloader.Download(self.rss_addr))
}

func (self *TumblrDownloader) get_image_list() [][]string{
    finder, _ := regexp.Compile(`img *?src="(.*?)"`)
    text := self.download_rss()
    rslts := finder.FindAllStringSubmatch(text, -1)
    return rslts
}

func (self *TumblrDownloader) check_rss(){
    for {
        image_list := self.get_image_list()
        log.Printf("Getting Rss %s: getted %d", self.rss_addr, len(image_list))
        //len(rslts)
        for _, v := range image_list{
            url := v[1]
            go self.ProcessUrl(url)
        }
        time.Sleep(time.Minute * Config.CheckInterval)
    }
}

func (self *TumblrDownloader) Start(){
    go self.check_rss()
    go self.AfterFinished()
    for {
        img_resource, ok := <- self.UrlChan
        if !ok {
            break
        }
        go Download_raw(img_resource, *self)
    }
}
