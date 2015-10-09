package database

import (
	//"fmt"
	"strconv"
	"strings"
)

type DBPostgres struct {
	DBase
}

func (db *DBPostgres) ReplaceParamsSymbol(sql *string) {
	num := 0
	new_sql := ""
	for _, c := range *sql {
		if c == '?' {
			num += 1
			new_sql += "$" + strconv.Itoa(num)
		} else {
			new_sql += string(c)
		}
	}
	*sql = new_sql
}

func (db *DBPostgres) QuoteIdentifier(name string) string {
	end := strings.IndexRune(name, 0)
	if end > -1 {
		name = name[:end]
	}
	return `"` + strings.Replace(name, `"`, `""`, -1) + `"`
}

func (db *DBPostgres) HasReturningId() bool {
	//Nota: ver el caso que la tabla no tenga id serial
	return true
}
