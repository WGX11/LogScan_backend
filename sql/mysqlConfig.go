package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MySQLConfig struct {
	URL string
}

var SqlConfig *MySQLConfig
var Db *sql.DB

func NewMySQLConfig() *MySQLConfig {
	return &MySQLConfig{
		URL: "wgx:wanguangxi.1@tcp(bj-cynosdbmysql-grp-279qdsqc.sql.tencentcdb.com:28647)" +
			"/logscan?charset=utf8mb4&parseTime=True&loc=Local",
	}
}

func InitDB() {
	var err error
	SqlConfig = NewMySQLConfig()
	Db, err = sql.Open("mysql", SqlConfig.URL)
	if err != nil {
		log.Fatalf("Failed to open mysql: %v", err)
	}
	if err = Db.Ping(); err != nil {
		log.Fatalf("Failed to ping mysql: %v", err)
	}
}

func CloseDB() {
	if Db != nil {
		Db.Close()
	}
}
