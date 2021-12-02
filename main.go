package main

import (
	"youtube/my"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	my.Migrate()
}
