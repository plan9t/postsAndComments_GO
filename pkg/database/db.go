package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectToDB создает подключение к базе данных и возвращает экземпляр *gorm.DB.
func ConnectToDB() (*gorm.DB, error) {
	dsn := "host=localhost user=plan9t dbname=ozon sslmode=disable password=plan9t"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fmt.Println("DB CONNECTED")
	if err != nil {
		return nil, err
	}

	return db, nil
}
