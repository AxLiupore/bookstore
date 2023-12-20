package dao

import "fmt"

func migration() {
	err := Db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate()
	if err != nil {
		fmt.Printf("err", err)
	}
	return
}
