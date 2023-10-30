package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

func InitDB() {
	DB, err = gorm.Open(
		postgres.Open(fmt.Sprintf("host=localhost user=postgres password=pass dbname=postgres port=5432 sslmode=disable")),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)},
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return DB
}
