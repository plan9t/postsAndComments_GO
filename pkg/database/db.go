package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Для docker-compose использовать "host.docker.internal" вместо "localhost" !!! Но надо потестить.

// ConnectToDB создает подключение к базе данных и возвращает экземпляр *gorm.DB.
func ConnectToDB() (*gorm.DB, error) {
	dsn := "host=host.docker.internal user=plan9t dbname=ozon sslmode=disable password=plan9t"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
