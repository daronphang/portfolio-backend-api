package dblayer

import "database/sql"

type DBLayer struct {
	*sql.DB
}

func DBInstance(instance string) (*DBLayer, bool) {
	switch instance {
	case "MSSQL":
		return MSSQLDB, true
	case "MYSQL":
		return MYSQLDB, true
	default:
		return nil, false
	}
}
