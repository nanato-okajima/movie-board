package main

import (
	"fmt"
	"youtube/my"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var envpath = "./.env"

func main() {

	loadEnv(envpath)
	my.Migrate()
}

func loadEnv(envpath string) {
	err := godotenv.Load(envpath)
	if err != nil {
		fmt.Println("error")
	}
}
