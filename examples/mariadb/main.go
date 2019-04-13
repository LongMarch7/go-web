package main

import(
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main(){

    db, err := sql.Open("mysql", "root:123456@tcp(localhost:13306)/test?charset=utf8")
    if err != nil {
        fmt.Println(err)
        return
    }
    db.SetMaxOpenConns(2000)
    db.SetMaxIdleConns(1000)
    defer db.Close()

    rows,err := db.Query("select * from test")
    defer rows.Close()

    columns, _ := rows.Columns()
    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    for rows.Next() {
        //将行数据保存到record字典
        err = rows.Scan(scanArgs...)
        for i, col := range values {
            if col != nil {
                fmt.Println("key", columns[i], "data", string(col.([]byte)))
            }
        }
    }
}