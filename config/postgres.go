package config

import (
	"fmt"

	"github.com/jinzhu/gorm"

	// gorm postgres url
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB is database instance
var DB *gorm.DB

func init() {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Env.DbHost, Env.DbPort, Env.DbUser, Env.DbName, Env.DbPassword,
	)

	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		panic(err)
	}

	DB = db

	fmt.Println("You connected to your database.")
}
