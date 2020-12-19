// This is example intend to call `UPPER` function which builtin MySQL

package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User Model
type SqlExpressUser struct {
	Name     string
	Location string
}

func main() {
	dsn := "root:redhat@tcp(localhost:3306)/testing?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to MySQL server.")
	}

	db.AutoMigrate(&SqlExpressUser{})

	db.Model(&SqlExpressUser{}).Create(map[string]interface{}{
		"Name": "JinZhu",
		"Location": clause.Expr{
			SQL:  "UPPER(?)",
			Vars: []interface{}{"HuNan ChangSha"},
		},
	})
}
