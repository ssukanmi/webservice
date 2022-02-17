package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/ssukanmi/webservice/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabaseConnection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to create a connection to database")
	}
	db.AutoMigrate(&entity.User{})
	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Println("Failed to close connection from database")
	}
	dbSQL.Close()
}