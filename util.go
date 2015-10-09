package simplify

import (
	"fmt"
	"github.com/jdroguett/simplify/database"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

//FIXME no se usa
func getColumsWithoutId(dbase database.DBInter, model interface{}) (columns []string) {
	columns = getColums(dbase, model)
	var strSlice sort.StringSlice = columns
	strSlice.Sort()
	pos := sort.SearchStrings(columns, "\"id\"")
	return append(columns[:pos], columns[pos+1:]...)
}

//Return ?, ?, ?
func getParams(numFields int) string {
	p := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		p[i] = "?"
	}
	return strings.Join(p, ", ")
}

func getParamsUpdate(columns []string) string {
	p := make([]string, len(columns))
	for i := 0; i < len(columns); i++ {
		p[i] = columns[i] + "=?"
	}
	return strings.Join(p, ", ")
}

func getColums(dbase database.DBInter, model interface{}) (columns []string) {
	ind := reflect.Indirect(reflect.ValueOf(model))
	typeOfM := ind.Type()
	for i := 0; i < ind.NumField(); i++ {
		columns = append(columns, dbase.QuoteIdentifier(strings.ToLower(typeOfM.Field(i).Name)))
	}
	return columns
}

func getTableName(model interface{}) string {
	ind := reflect.Indirect(reflect.ValueOf(model))
	return strings.ToLower(ind.Type().Name())
}

func getModelId(model interface{}) string {
	return "Id"
}

func getColumsAndData(dbase database.DBInter, model interface{}) (columns []string, data []interface{}) {
	ind := reflect.Indirect(reflect.ValueOf(model))
	typeOfM := ind.Type()
	for i := 0; i < ind.NumField(); i++ {
		if typeOfM.Field(i).Name != "Id" {
			columns = append(columns, dbase.QuoteIdentifier(strings.ToLower(typeOfM.Field(i).Name)))
			data = append(data, ind.Field(i).Interface())
		}
	}
	return columns, data
}

func getValueId(model interface{}) (Id int) {
	ind := reflect.Indirect(reflect.ValueOf(model))
	return ind.FieldByName(getModelId(model)).Interface().(int)
}

func getData(model interface{}) {
	ind := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < ind.NumField(); i++ {
		a := ind.Field(i)
		fmt.Println(a)
	}

}

func setModel(model interface{}, columns []string, values []interface{}) {
	mod := reflect.Indirect(reflect.ValueOf(model))
	for i, name := range columns {
		field := mod.FieldByName(strings.Title(name))
		val := reflect.Indirect(reflect.ValueOf(values[i])).Interface()

		switch v := val.(type) {
		case int64:
			field.SetInt(v)
		case []byte: //mysql parece que devuelve todo como []byte
			if field.Kind() == reflect.Int {
				n, _ := strconv.Atoi(string(v))
				field.SetInt(int64(n))
			} else {
				field.SetString(string(v))
			}
		case string:
			field.SetString(v)
		}
	}
}
