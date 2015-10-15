package simplify

import (
	"database/sql"
	"fmt"
	"github.com/jdroguett/simplify/database"
	"reflect"
	"strings"
)

type Model struct {
	Db       *sql.DB
	DBase    database.DBInter
	WhereStr string
	Args     []interface{}
	OrderStr string
}

func Open(driverName, dataSourceName string) (m *Model, err error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	var dBase database.DBInter
	switch driverName {
	case "postgres":
		dBase = new(database.DBPostgres)
	case "mysql":
		dBase = new(database.DBMysql)
	case "sqlite3":
		dBase = new(database.DBSqlite3)
	}
	return &Model{Db: db, DBase: dBase}, nil

}

func (m *Model) Close() error {
	err := m.Db.Close()
	if err != nil {
		return err
	}
	return nil
}

// Example struct:
//
// type User struct {
//    Id    int
//    Age   int
//    Name  string
//    Email string
// }
//
// Use:
//
//    var users []User
//    err = orm.Query( &users, "SELECT * FROM \"user\" order by name asc")
func (m *Model) Query(models interface{}, query string) (err error) {
	ind := reflect.Indirect(reflect.ValueOf(models))
	elem := ind.Type().Elem()

	Log.Println(query)
	rows, err := m.Db.Query(query)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	defer rows.Close()

	values := make([]interface{}, len(columns))
	for i := range values {
		var val interface{}
		values[i] = &val
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return err
		}
		model := reflect.New(elem)
		setModel(model.Interface(), columns, values)
		ind.Set(reflect.Append(ind, reflect.Indirect(model)))
	}
	return nil
}

// Use:
//
//     user := User{Name: "jean", Email: "x@x.com", Age: 40}
//     orm.Insert(&user)
// Result:
//     user == User{Id: 3223, Name: "jean", Email: "x@x.com", Age: 40}
func (m *Model) Insert(model interface{}) (err error) {

	ind := reflect.Indirect(reflect.ValueOf(model))
	columns, d := getColumsAndData(m.DBase, model)
	var id int64

	if m.DBase.HasReturningId() {
		sql := fmt.Sprintf("INSERT INTO %v(%v) VALUES (%v) RETURNING %v",
			m.DBase.QuoteIdentifier(getTableName(model)),
			strings.Join(columns, ", "),
			getParams(len(columns)),
			m.DBase.QuoteIdentifier(strings.ToLower(getModelId(model))))
		m.DBase.ReplaceParamsSymbol(&sql)
		Log.Println(sql, d)
		err = m.Db.QueryRow(sql, d...).Scan(&id)
		if err != nil {
			return err
		}
	} else {
		sql := fmt.Sprintf("INSERT INTO %v(%v) VALUES (%v)",
			m.DBase.QuoteIdentifier(getTableName(model)),
			strings.Join(columns, ", "),
			getParams(len(columns)))
		m.DBase.ReplaceParamsSymbol(&sql)
		Log.Println(sql, d)
		stmt, err := m.Db.Prepare(sql)
		if err != nil {
			return err
		}
		defer stmt.Close()
		res, err := stmt.Exec(d...)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}
	}

	fieldId := ind.FieldByName(getModelId(model))
	var v interface{} = int64(id)
	fieldId.Set(reflect.ValueOf(v))
	return nil
}

// Use:
//
//     user := User{Id: 32, Name: "jean", Email: "x@x.com", Age: 40}
//     orm.Update(user)
func (m *Model) Update(model interface{}) (err error) {
	columns, data := getColumsAndData(m.DBase, model)

	sql := fmt.Sprintf("UPDATE %v SET %v WHERE %v=%v",
		m.DBase.QuoteIdentifier(getTableName(model)),
		getParamsUpdate(columns),
		m.DBase.QuoteIdentifier(strings.ToLower(getModelId(model))),
		"?")

	data = append(data, getValueId(model))
	m.DBase.ReplaceParamsSymbol(&sql)
	Log.Println(sql, data)
	stmt, err := m.Db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(data...)
	if err != nil {
		return err
	}

	return nil
}

// Insert or Update
//
// Use:
//     user := User{Id: 32, Name: "jean", Email: "x@x.com", Age: 40}
//     orm.save(&user)
//
// id is nil or 0 => insert, id is not nil => update
func (m *Model) Save(model interface{}) (err error) {
	id := getValueId(model)
	if id == 0 {
		err = m.Insert(model)
		if err != nil {
			return err
		}
	} else {
		err = m.Update(model)
		if err != nil {
			return err
		}
	}
	return nil
}

// Use:
//
//     user := User{Id: 32}
//     orm.Delete(user)
func (m *Model) Delete(model interface{}) (err error) {
	var data []interface{}
	sql := fmt.Sprintf("DELETE FROM %v WHERE %v=%v",
		m.DBase.QuoteIdentifier(getTableName(model)),
		m.DBase.QuoteIdentifier(strings.ToLower(getModelId(model))),
		"?")
	data = append(data, getValueId(model))
	m.DBase.ReplaceParamsSymbol(&sql)
	Log.Println(sql, data)
	stmt, err := m.Db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(data...)
	if err != nil {
		return err
	}
	return nil
}

// Use:
//    var user User
//	  err = orm.Where("email = $1", "xyz@x.com").Order("id desc").First(&user)
func (m *Model) Where(where string, args ...interface{}) *Model {
	m.WhereStr = where
	m.Args = args
	return m
}

// Use:
//    var user User
//	  err = orm.Where("email = $1", "xyz@x.com").Order("id desc").First(&user)
func (m *Model) Order(order string) *Model {
	m.OrderStr = order
	return m
}

// Use:
//    var user User
//	  err = orm.Where("email = $1", "xyz@x.com").Order("id desc").First(&user)
func (m *Model) First(model interface{}) (err error) {
	sql := fmt.Sprintf("SELECT * FROM %v WHERE (%v)",
		m.DBase.QuoteIdentifier(getTableName(model)),
		m.WhereStr)
	if m.OrderStr != "" {
		sql = fmt.Sprintf("%v ORDER BY %v", sql, m.OrderStr)
	}
	sql += " limit 1"

	m.DBase.ReplaceParamsSymbol(&sql)
	Log.Println(sql, m.Args)

	stmt, err := m.Db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(m.Args...)
	if err != nil {
		return err
	}
	defer res.Close()

	if res.Next() {
		columns, err := res.Columns()
		if err != nil {
			return err
		}
		values := make([]interface{}, len(columns))
		for i := range values {
			var val interface{}
			values[i] = &val
		}

		err = res.Scan(values...)
		if err != nil {
			return err
		}
		setModel(model, columns, values)
	}
	return nil
}
