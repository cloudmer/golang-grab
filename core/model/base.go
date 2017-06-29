package model

import (
	"database/sql"
	"xmn/core/config"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
	"strconv"
)

var DB *sql.DB
var err error

func init()  {
	fmt.Println("in connet mysql ... ")
	db_type := "mysql"
	db_user := config.Read("mysql", "user")
	db_pass := config.Read("mysql", "pass")
	db_host := config.Read("mysql", "host")
	db_name := config.Read("mysql", "name")
	db_port := config.Read("mysql", "port")
	db_max_open := config.Read("mysql", "MaxOpen")
	max_open, _ := strconv.Atoi(db_max_open)
	db_max_idle := config.Read("mysql", "MaxIdle")
	max_idle, _ := strconv.Atoi(db_max_idle)
	DB, err = sql.Open(db_type, db_user+":"+db_pass+"@tcp("+db_host+":"+db_port+")/"+db_name+"?charset=utf8")
	DB.SetMaxOpenConns(max_open)
	DB.SetMaxIdleConns(max_idle)
	err = DB.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}