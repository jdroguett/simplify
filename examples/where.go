package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jdroguett/simplify"
	_ "github.com/lib/pq"
)

// User model: create table "user"(id serial, name varchar, email varchar, created_at timestamp, updated_at timestamp);
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func main() {
	//logger
	simplify.Debug = true

	sim, err := simplify.Open("postgres", "user=basego dbname=basego password=basego sslmode=disable")
	checkErr(err)
	defer sim.Close()

	//query
	var users []User
	err = sim.Where("name = ?", "jean1").All(&users)
	checkErr(err)

	for i, user := range users {
		fmt.Println("user: ", i, user)
	}

	jsonUsers, err := json.MarshalIndent(users, "", "")
	checkErr(err)
	fmt.Println("json users: ", string(jsonUsers))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
