package main

import (
	"fmt"
	"time"

	"github.com/jdroguett/simplify"
	_ "github.com/lib/pq"
)

// User model: create table "user"(id serial, name varchar, email varchar, created_at timestamp, updated_at timestamp);
type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	//logger
	simplify.Debug = true

	sim, err := simplify.Open("postgres", "user=basego dbname=basego password=basego sslmode=disable")
	checkErr(err)
	defer sim.Close()

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
