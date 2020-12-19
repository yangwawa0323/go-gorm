package main

import (
	"database/sql"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User struct
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
	// Use gorm open the mysql driver open connection with gorm.Config
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect MySQL server. Please contact administrator.")
	}

	// Migrate schema
	db.AutoMigrate(&User{})

	// Create a user
	now := time.Now()
	user := User{Name: "Jinzhu", Age: 18, Birthday: &now}
	_ = db.Create(&user)
	time.Sleep(30 * time.Second)
	db.Delete(&user)
}
