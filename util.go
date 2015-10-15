package simplify

import (
	"fmt"
	"github.com/jdroguett/simplify/database"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
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
			columns = append(columns, dbase.QuoteIdentifier(columnName(typeOfM.Field(i).Name)))
			field := ind.Field(i).Interface()
			switch field.(type) {
			case time.Time:
				data = append(data, field.(time.Time).UTC())
			default:
				data = append(data, field)
			}
		}
	}
	return columns, data
}

func getValueId(model interface{}) (Id int64) {
	ind := reflect.Indirect(reflect.ValueOf(model))
	return ind.FieldByName(getModelId(model)).Interface().(int64)
}

//no se usa
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
		field := mod.FieldByName(fieldName(name))
		if field.IsValid() {
			val := reflect.Indirect(reflect.ValueOf(values[i])).Interface()
			switch v := val.(type) {
			case int64:
				field.SetInt(v)
			case []byte: //mysql parece que devuelve todo como []byte
				switch field.Kind() {
				case reflect.Int64:
					n, _ := strconv.Atoi(string(v))
					field.SetInt(int64(n))
				case reflect.Struct:
					t, _ := time.Parse("2006-01-02 15:04:05", string(v))
					field.Set(reflect.ValueOf(t))
				default:
					field.SetString(string(v))
				}
			case string:
				field.SetString(v)
			case time.Time:
				field.Set(reflect.ValueOf(v))
			} //end switch
		}
	}
}

//Name => name
//CreatedAt => created_at
func columnName(name string) (new_name string) {
	new_name = ""
	for i, runeValue := range name {
		if i > 0 && unicode.IsUpper(runeValue) {
			new_name += "_"
		}
		new_name += string(unicode.ToLower(runeValue))
	}
	return new_name
}

//name => Name
//created_at => CreatedAt
func fieldName(name string) (new_name string) {
	new_name = ""
	b := false
	for i, runeValue := range name {
		if i == 0 || b {
			new_name += string(unicode.ToUpper(runeValue))
			b = false
		} else if runeValue == '_' {
			b = true
		} else {
			new_name += string(runeValue)
		}
	}
	return new_name
}
