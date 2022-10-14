package dblayer

import (
	"database/sql"
	"portfolio_golang/src/config"
	"portfolio_golang/src/zaplog"
	"time"

	"github.com/go-sql-driver/mysql"
)

var MYSQLDB = &DBLayer{}

func InitMYSQL(appCfg *config.EnvConf) (*DBLayer, error) {
	// DSN 
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname?timeout=90s&collation=utf8mb4_unicode_ci

	// for local sockets: Addr = /var/run/mysqld/mysqld.sock, Net = unix
	// ipv6 must be enclosed with []

	cfg := &mysql.Config{
		User: appCfg.MYSQLUSERNAME,
		Passwd: appCfg.MYSQLPASSWORD,
		Net: appCfg.MYSQLPROTOCOL,
		Addr: appCfg.MYSQLADDRESS, 
		DBName: "portfolio",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	// important settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// initializing global variable
	MYSQLDB.DB = db

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	zaplog.Logger.Info("MYSQL connection established")
	return MYSQLDB, nil
}
