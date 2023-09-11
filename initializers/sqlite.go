package initializers

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lightsaid/blogs/config"
	_ "github.com/mattn/go-sqlite3"
)

// InitSQLite 连接SQLite并设置基础参数
func InitSQLite() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", config.AppConf.Database.Source)
	if err != nil {
		return db, err
	}

	db.SetMaxOpenConns(config.AppConf.Database.MaxOpenConns)
	db.SetMaxIdleConns(config.AppConf.Database.MaxIdleConns)
	maxIdleTime := config.ParseDuration(config.AppConf.Database.MaxIdleTime, 5*time.Minute)
	db.SetConnMaxIdleTime(maxIdleTime)

	return db, nil
}
