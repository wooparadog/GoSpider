package main

type Recorder interface {
    MarkAsFinished(url string)
    HasFinished(url string) bool
}

func CreateRecorders() Recorder{
    var recorder Recorder
    if Config.Recorder == "sqlite"{
        sqlite_recorder := CreateSqliteRecorder(SQLITE_FILE)
        recorder = &sqlite_recorder
    }else if Config.Recorder == "redis"{
        redis_recorder := CreateRedisRecord(REDIS_NETWORK, REDIS_LOCATION)
        recorder = &redis_recorder
    }
    return recorder
}
