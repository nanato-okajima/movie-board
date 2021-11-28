package my

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", host, user, password, dbname, port)

	db, er := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if er != nil {
		fmt.Println(er)
		return
	}

	db.AutoMigrate(&User{}, &Group{}, &Post{}, &Comment{})
}
