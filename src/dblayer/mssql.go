package dblayer

import (
	"database/sql"
	"fmt"
	"net/url"
	"portfolio_golang/src/config"
	"portfolio_golang/src/zaplog"

	_ "github.com/denisenkom/go-mssqldb"
)

var MSSQLDB = &DBLayer{}

func InitMSSQL(appCfg *config.EnvConf) (*DBLayer, error) {
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(appCfg.MSSQLUSERNAME, appCfg.MSSQLPASSWORD),
		Host:   fmt.Sprintf("%s:%d", appCfg.MSSQLSERVER, appCfg.MSSQLPORT),
	}

	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		return nil, err
	}
	// initializing global variable
	MSSQLDB.DB = db

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	zaplog.Logger.Info("MSSQL connection established")
	return MSSQLDB, nil
}
