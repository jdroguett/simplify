package database

import (
//"strings"
)

type DBInter interface {
	ReplaceParamsSymbol(*string)
	QuoteIdentifier(name string) string
	HasReturningId() bool
}

type DBase struct {
}

func (db *DBase) ReplaceParamsSymbol(query *string) {
	/* implement: default use ? */
}

func (db *DBase) QuoteIdentifier(name string) string {
	return name
}

func (db *DBase) HasReturningId() bool {
	//default
	return false
}
