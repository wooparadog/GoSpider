package main

import(
    "log"
    "sync"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type sqlite_row struct{
    url string
}

type SqliteRecorder struct{
    conn *sql.DB
    mutex sync.Mutex
}

func CreateSqliteRecorder(location string) SqliteRecorder{
    conn, err := sql.Open("sqlite3", location)
    if err != nil {
        log.Fatal("Sqlite connection fail!\nFile: ", SQLITE_FILE)
    }
    _, err = conn.Exec("create table if not exists download_record(url varchar(512) primary key) ")
    if err != nil{
        log.Fatal("Sqlite Recorder Table download_record does not exists and cannot be created")
    }
    log.Println("Sqlite Recorder Connected : ", location)
    return SqliteRecorder{conn: conn}
}

func (self *SqliteRecorder) MarkAsFinished(url string){
    self.mutex.Lock()
    self.conn.Exec("insert into download_record values(?)", url)
    self.mutex.Unlock()
}

func (self *SqliteRecorder) HasFinished(url string) bool{
    self.mutex.Lock()
    row := self.conn.QueryRow("select url from download_record where url=?", url)
    self.mutex.Unlock()
    sqlite_record := sqlite_row{}
    if err:=row.Scan(&sqlite_record.url); err!=nil{
        return false
    }
    return true
}
