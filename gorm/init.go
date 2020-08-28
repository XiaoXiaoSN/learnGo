package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	var err error

	user := "test"
	pw := "test"
	address := "localhost:33333"
	dbName := "test"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, pw, address, dbName)

	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(true)

	migration()
}

func migration() {
	db.AutoMigrate(&User{})
	fmt.Println("migration done.")
}
