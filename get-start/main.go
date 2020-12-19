package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	dsn := "root:redhat@tcp(127.0.0.1:3306)/testing?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&Product{})

	db.Create(&Product{
		Code:  "E42",
		Price: 100,
	})

	// Read
	var product Product
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "E42") // find product with code D42

	fmt.Printf("find E42 :%v\n", &product)

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)

	fmt.Printf("update E42 :%v\n", &product)
	// Update - update multiple fields
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.Delete(&product, 3)
	db.First(&product, 3)
	fmt.Printf("find E42 :%v\n", &product)
}
