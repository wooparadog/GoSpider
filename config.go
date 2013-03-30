package main

import (
    "os"
    "log"
    "time"
    "path"
    "encoding/json"
    "io/ioutil"
)
const CONCURENT_DOWNLOADS = 4

const REDIS_NETWORK = "unix"
const REDIS_LOCATION = "/tmp/redis.sock"

const SQLITE_FILE = "data.sqlite3"

const SET_NAME = "url_set"

type TumblrSource struct{
    Name string
    Suffix string
    Url string
}

type config_struct struct{
    Proxy string
    CheckInterval time.Duration
    Recorder string
    UseProxy bool
    TumblrSources []TumblrSource
    Timeout int64
}

var Config config_struct

func ParseConfig(){
    dir, _ := os.Getwd()
    config_file := path.Join(dir, "config.json")
    config, _ := os.Open(config_file)
    raw_config, err := ioutil.ReadAll(config)
    if err == nil{
        err = json.Unmarshal(raw_config, &Config)
        log.Printf("Proxy:\t\t %s", Config.Proxy)
        log.Printf("Check Interval:\t %s", Config.CheckInterval)
        log.Printf("Tumblr Sources:\t %s", Config.TumblrSources)
        return
    }
    if err != nil{
        log.Fatal("Failed To Load Config File")
    }
}

