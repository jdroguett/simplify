# Simplify
Simplifying access to data

Drivers for Go's sql package which support database/sql

Tested with:
* PostgreSQL: github.com/lib/pq
* Mysql: github.com/go-sql-driver/mysql
* SQLite3: github.com/mattn/go-sqlite3

## Conventions

Column name (DB) | Field name (struct)
------------- | -------------
Id  | id
Name  | name
CreatedAt  | created_at
UpdatedAt  | updated_at
Etc | etc


## Example
```go
package main

import (
	"fmt"
	"github.com/jdroguett/simplify"
	_ "github.com/lib/pq"
	"time"
)

/* create table "user"(id serial, name varchar, email varchar, created_at timestamp, updated_at timestamp); */
type User struct {
	Id        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	//logger
	simplify.Debug = true

	sim, err := simplify.Open("postgres", "user=basego dbname=basego sslmode=disable")
	checkErr(err)
	defer sim.Close()

	//insert
	user := User{Name: "Jean", Email: "x@x.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = sim.Insert(&user)
	checkErr(err)
	fmt.Println("user: ", user)

	//select one
	var user2 User
	err = sim.Where("name = ?", "Jean").Order("id desc").First(&user2)
	checkErr(err)
	fmt.Println("user2: ", user2)

	//update
	user = User{Id: 1, Name: "Jean update", Email: "update@x.com", UpdatedAt: time.Now()}
	err = sim.Update(user)
	checkErr(err)

	//delete
	user_del := User{Id: 2}
	err = sim.Delete(user_del)
	checkErr(err)

	//insert or update
	user = User{Name: "user new", Email: "xyz@x.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = sim.Save(&user)
	checkErr(err)

	//query
	var users []User
	err = sim.Query(&users, "SELECT * FROM \"user\" ORDER BY name ASC")
	fmt.Println("users: ", users)
	checkErr(err)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

```
