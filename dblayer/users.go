package dblayer

import (
	"fmt"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	_ "time"
)

func GetAllUsers() []string {
	fmt.Println("getting all users")
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
