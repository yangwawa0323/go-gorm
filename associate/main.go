package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// AssUser struct
type AssUser struct {
	gorm.Model
	Name       string
	CreditCard AssCreditCard
}

// AssCreditCard struct
type AssCreditCard struct {
	gorm.Model
	Number    string
	AssUserID uint
}

func main() {
	dsn := "root:redhat@tcp(localhost:3306)/testing?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to MySQL server.")
	}

	db.AutoMigrate(&AssUser{}, &AssCreditCard{})

	db.Create(&AssUser{
		Name:       "jinzhu",
		CreditCard: AssCreditCard{Number: "411111111111"},
	})

}
