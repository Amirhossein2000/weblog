package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"weblog/schema"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("sqlite3", "./weblog.db")
	if err != nil {
		log.Panic("Database err:", err)
	}
	err = DB.AutoMigrate(&schema.User{}, &schema.Article{}, &schema.Comment{}).Error
	if err != nil {
		log.Panic("Database err:", err)
	}
}
