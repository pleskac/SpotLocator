package dblayer

import (
	"fmt"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"strconv"
)

func GetLatestSpotId() int {
	spotStr := getValue(lATEST_SPOT)
	spotId, err := strconv.Atoi(spotStr)

	if err != nil {
		fmt.Println(err)
		return -1
	}

	return spotId
}

func SaveLatestSpotId(id int) {
	idStr := fmt.Sprintf("%d", id)
	saveValue(lATEST_SPOT, idStr)
}

func getValue(idType string) string {
	db := Connect()
	defer db.Close()

	rows, _, err := db.Query("SELECT * FROM kvp WHERE k = '" + idType + "'")
	if err != nil {
		panic(err)
	}

	if len(rows) < 1 {
		return "-1"
	}

	return rows[0].Str(0)
}

func saveValue(key string, newValue string) {
	db := Connect()
	defer db.Close()

	//Delete that row, if it exists
	stmt, err := db.Prepare("DELETE FROM kvp WHERE k = '" + key + "'")
	_, err = stmt.Run()
	if err != nil {
		//row did not exist, we'll just add it later
	}

	//Insert new row
	stmt, err = db.Prepare("INSERT INTO kvp (k, v) VALUES (?, ?)")
	_, err = stmt.Run(key, newValue)
	if err != nil {
		panic(err)
	}
}
