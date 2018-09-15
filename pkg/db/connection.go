package db

import (
	"time"

	"github.com/gocraft/dbr"

	"kubesphere.io/devops/pkg/config"
)

var defaultEventReceiver = EventReceiver{}

func OpenDatabase(cfg config.MysqlConfig) (*Database, error) {
	// https://github.com/go-sql-driver/mysql/issues/9
	conn, err := dbr.Open("mysql", cfg.GetUrl()+"?parseTime=1&multiStatements=1&charset=utf8mb4&collation=utf8mb4_unicode_ci", &defaultEventReceiver)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(100)
	conn.SetMaxOpenConns(100)
	conn.SetConnMaxLifetime(10 * time.Second)
	return &Database{
		Session: conn.NewSession(nil),
	}, nil
}
