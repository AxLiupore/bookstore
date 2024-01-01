package dao

import (
	"bookstore/model"
	"fmt"
)

func migration() {
	err := Db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		&model.Address{},
		&model.Admin{},
		&model.BasePage{},
		&model.Carousel{},
		&model.Cart{},
		&model.Category{},
		&model.Favorite{},
		&model.Notice{},
		&model.Order{},
		&model.Product{},
		&model.ProductImg{},
		&model.User{})
	if err != nil {
		fmt.Printf("err", err)
	}
	return
}
