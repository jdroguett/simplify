package database

import ()

type DBMysql struct {
	DBase
}

func (db *DBMysql) QuoteIdentifier(name string) string {
	return "`" + name + "`"
}
