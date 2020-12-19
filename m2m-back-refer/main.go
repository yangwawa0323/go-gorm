package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MmbrUser back-reference
type MmbrUser struct {
	gorm.Model
	Name      string
	Location  string
	Email     string
	Languages []*MmbrLanguage `gorm:"many2many:mmbr_user_languages"`
}

// MmbrLanuguage Back reference
type MmbrLanguage struct {
	gorm.Model
	Name  string
	Users []*MmbrUser `gorm:"many2many:mmbr_user_languages"`
}

func main() {
	dsn := "root:redhat@tcp(localhost:3306)/testing?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Can not connect to MySQL server.")
	}
	db.Set("gorm:table_options", "DEFAULT CHARSET utf8mb4").AutoMigrate(&MmbrUser{}, &MmbrLanguage{})

	han := MmbrLanguage{Name: "汉语"}
	// han := MmLanguage{Name: "Han"}
	english := MmbrLanguage{Name: "English"}
	user := MmbrUser{Name: "Yangwawa",
		Email:     "Yangwawa0323@163.com",
		Languages: []*MmbrLanguage{&han, &english},
	}

	user2 := MmbrUser{
		Name:      "jinzhu",
		Email:     "fake@qq.com",
		Languages: []*MmbrLanguage{&han},
	}
	users := []MmbrUser{user, user2}
	db.Create(&users)
}
