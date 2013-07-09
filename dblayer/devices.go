package dblayer

import (
	_ "fmt"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	_ "time"
)

func GetDevices(user string) []Device {
	var devices []Device
	db := Connect()
	defer db.Close()

	//Get the current trip, if it exists
	rows, _, err := db.Query("SELECT * FROM devices WHERE user_email = " + user)
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		devices = append(devices, Device{row.Str(0), row.Str(1), row.Str(2), row.Str(3)})
	}

	return devices
}
