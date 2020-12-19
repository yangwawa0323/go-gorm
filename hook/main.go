package main

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User Model
type User struct {
	ID           uint
	UUID         uuid.UUID
	Name         string
	Email        *string
	Age          uint8
	Role         string
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivedAt    sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Hook before create a record
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()

	if u.Role == "admin" {
		return errors.New("invalid role")
	}
	return
}

func main() {
	dsn := "root:redhat@tcp(localhost:3306)/testing?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to MySQL server.")
	}

	// Any struct modification need migrate into database.
	db.AutoMigrate(&User{})

	// Select fileds for insertion, only seleted field can insert into table.
	now := time.Now()
	email := "JinZhu@github.com"
	userPtr := &User{Name: "Jinzhu", Age: 18, Birthday: &now, Email: &email}
	db.Create(userPtr)

}
