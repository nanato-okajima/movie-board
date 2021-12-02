package my

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"youtube/config"
)

var DB *gorm.DB

func init() {
	conf := config.DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", conf.Host, conf.User, conf.Password, conf.DBName, conf.Port)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(dsn + "database can't connect")
	}
	DB.AutoMigrate(&User{}, &Post{}, &Group{}, &Comment{})
}
