package go_spider

import(
    "log"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type sqlite_row struct{
    url string
}

type SqliteRecorder struct{
    conn chan *sql.DB
}

var ConnectionPool chan *sql.DB

func CreateSqliteRecorder(location string) SqliteRecorder{
    conn, err := sql.Open("sqlite3", location)
    if err != nil {
        log.Fatal("Sqlite connection fail!\nFile: ", SQLITE_FILE)
    }
    _, err = conn.Exec("create table if not exists download_record(url varchar(512) primary key) ")
    if err != nil{
        log.Fatal("Sqlite Recorder Table download_record does not exists and cannot be created")
    }
    ConnectionPool = make(chan *sql.DB, 1)
    ConnectionPool <- conn
    log.Println("Sqlite Recorder Connected : ", location)
    return SqliteRecorder{conn: ConnectionPool}
}

func (self *SqliteRecorder) Execute(sql string, params ...interface{}) (sql.Result, error){
    conn := <- self.conn
    defer func(){self.conn <- conn}()
    return conn.Exec(sql, params...)
}

func (self *SqliteRecorder) QueryRow(sql string, params ...interface{}) *sql.Row{
    conn := <- self.conn
    defer func(){self.conn <- conn}()
    return conn.QueryRow(sql, params...)
}

func (self *SqliteRecorder) MarkAsFinished(url string){
    self.Execute("insert into download_record values(?)", url)
}

func (self *SqliteRecorder) HasFinished(url string) bool{
    row := self.QueryRow("select url from download_record where url=?", url)
    sqlite_record := sqlite_row{}
    if err:=row.Scan(&sqlite_record.url); err!=nil{
        return false
    }
    return true
}
