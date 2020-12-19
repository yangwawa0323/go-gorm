package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MmUser struct {
	gorm.Model
	Name     string
	Location string
	Email    string
	Language []MmLanguage `gorm:"many2many:user_languages"`
}

type MmLanguage struct {
	gorm.Model
	Name string
	// Users []*MmUser `gorm:"many2many:user_languages"`
}

func main() {
	dsn := "root:redhat@tcp(localhost:3306)/testing?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Can not connect to MySQL server.")
	}
	db.Set("gorm:table_options", "DEFAULT CHARSET utf8mb4").AutoMigrate(&MmUser{}, &MmLanguage{})

	han := MmLanguage{Name: "汉语"}
	// han := MmLanguage{Name: "Han"}
	english := MmLanguage{Name: "English"}
	user := MmUser{Name: "Yangwawa",
		Email:    "Yangwawa0323@163.com",
		Language: []MmLanguage{han, english},
	}

	db.Create(&user)
}
