package main

import(
    "redis"
    "os"
    "fmt"
    "sync"
)

type RedisRecorder struct{
    conn *redis.Redis 
    mutex sync.Mutex
}

func CreateRedisRecord(network, location string) RedisRecorder{
    conn := redis.Redis{
        Network:network,
        Location:location,
    }
    if !conn.Connect() {
        fmt.Println("Redis connection fail!\nNetwork: ", network, " Location: ", location)
        os.Exit(1)
    }
    fmt.Println("Redis Recorder Connected")
    return RedisRecorder{conn: &conn}
}


func (self *RedisRecorder) MarkAsFinished(url string){
    self.mutex.Lock()
    self.conn.Sadd(SET_NAME, url)
    self.mutex.Unlock()
}

func (self *RedisRecorder) HasFinished(url string) bool{
    self.mutex.Lock()
    result := self.conn.Sismember(SET_NAME, url)
    self.mutex.Unlock()
    return result
}
