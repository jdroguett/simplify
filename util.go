package simplify

import (
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jdroguett/simplify/database"
)

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

func getColumns(dbase database.DBInter, model interface{}) (columns []string) {
	ind := reflect.Indirect(reflect.ValueOf(model))
	typeOfM := ind.Type()
	for i := 0; i < ind.NumField(); i++ {
		columns = append(columns, dbase.QuoteIdentifier(strings.ToLower(typeOfM.Field(i).Name)))
	}
	return columns
}

func getTableName(model interface{}) string {
	ind := reflect.Indirect(reflect.ValueOf(model))
	if ind.Kind() == reflect.Array || ind.Kind() == reflect.Slice {
		return strings.ToLower(ind.Type().Elem().Name())
	}
	return strings.ToLower(ind.Type().Name())
}

func getModelID(model interface{}) string {
	return "ID"
}

func getColumnsAndData(dbase database.DBInter, model interface{}) (columns []string, data []interface{}) {
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

func getValueID(model interface{}) (ID int64) {
	ind := reflect.Indirect(reflect.ValueOf(model))
	return ind.FieldByName(getModelID(model)).Interface().(int64)
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
					//t, _ := time.Parse("2006-01-02 15:04:05", string(v))
					//field.Set(reflect.ValueOf(t))
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

//ID => id
//Name => name
//CreatedAt => created_at
func columnName(name string) (newName string) {
	newName = ""
	for i, runeValue := range name {
		if i > 0 && unicode.IsUpper(runeValue) {
			newName += "_"
		}
		newName += string(unicode.ToLower(runeValue))
	}
	return newName
}

//id => ID
//name => Name
//created_at => CreatedAt
func fieldName(name string) string {
	if len(name) == 2 {
		return strings.ToUpper(name)
	}
	newName := ""
	b := false
	for i, runeValue := range name {
		if i == 0 || b {
			newName += string(unicode.ToUpper(runeValue))
			b = false
		} else if runeValue == '_' {
			b = true
		} else {
			newName += string(runeValue)
		}
	}
	return newName
}
