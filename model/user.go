package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string
	Password string
	Status   string
	Avatar   string
	Monty    string
}
