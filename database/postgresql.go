package database

import (
	//"fmt"
	"github.com/lib/pq"
	"strconv"
	//"strings"
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
	return pq.QuoteIdentifier(name)
}

func (db *DBPostgres) HasReturningId() bool {
	//Nota: ver el caso que la tabla no tenga id serial
	return true
}
