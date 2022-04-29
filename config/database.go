package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	sql "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/ssukanmi/webservice/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabaseConnection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	rootCert := os.Getenv("ROOT_CERT")

	gormLog, _ := os.Create("gorm.log")

	newLogger := logger.New(
		log.New(gormLog, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Warn,            // Log level
			IgnoreRecordNotFoundError: false,                  // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Disable color
		},
	)

	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(rootCert)
	if err != nil {
		fmt.Println("Failed to read root certificate path")
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		fmt.Println("Failed to append PEM.")
	}
	sql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs: rootCertPool,
	})
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=custom", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: "mysql",
		DSN:        dsn,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("Failed to create a connection to database")
	}
	db.AutoMigrate(&entity.User{}, &entity.UserImage{})
	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Println("Failed to close connection from database")
	}
	dbSQL.Close()
}
