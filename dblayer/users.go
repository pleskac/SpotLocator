package dblayer

import (
	_ "fmt"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	_ "time"
)

func GetAllUsers() []string {
	var users []string
	db := Connect()
	defer db.Close()

	//Get the current trip, if it exists
	rows, _, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		users = append(users, row.Str(0))
	}

	return users
}

func AddUser(user, username, display_name, password string) {
	db := Connect()
	defer db.Close()

	query := "INSERT INTO users (email, username, display_name, password) VALUES (?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	_, err = stmt.Run(user, username, display_name, password)
	if err != nil {
		panic(err)
	}
}

//authenticate functions
