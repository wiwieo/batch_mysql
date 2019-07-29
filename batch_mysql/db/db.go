package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"wiwieo/batch_mysql/constant"
)

var Conn *DB

func init() {
	Conn = ConnectToMySQL(constant.Config.Host, constant.Config.Port, constant.Config.DBName, constant.Config.User, constant.Config.Pwd)
}

type DB struct {
	*sqlx.DB
}

func ConnectToMySQL(host, port, dbName, user, pwd string) *DB{
	conn := sqlx.MustConnect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, host, port, dbName))
	conn.SetMaxOpenConns(50)
	conn.SetMaxIdleConns(10)
	return &DB{
		conn,
	}
}
