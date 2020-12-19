package main

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User Model
type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivedAt    sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func main() {
	dsn := "root:redhat@tcp(localhost:3306)/testing?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to MySQL server.")
	}

	// Select fileds for insertion, only seleted field can insert into table.
	now := time.Now()
	email := "JinZhu@github.com"
	userPtr := &User{Name: "Jinzhu", Age: 18, Birthday: &now, Email: &email}
	db.Select("Name", "Age", "Email").Create(userPtr)

	// Batch insertion
	var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
	db.Create(&users)

	for _, user := range users {
		fmt.Println("User id: ", user.ID) // 1,2,3
	}

}
