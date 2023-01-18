package config

import (
	"fmt"
	"golang/golang-skeleton/helper"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Set ip database connection
func SetUpDatabaseConnection() *gorm.DB {

	log.Println("[start load database conenction]")
	AppLoadEnv()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		helper.LogIfError(fmt.Errorf("failed to create connection to database, Error : %v", err))
	}

	//function to config database connection
	errConfigDB := configConnectionDB(db)
	helper.LogIfError(errConfigDB)

	// Auto Migrate
	if os.Getenv("MIGRATION") == "true" {
		migrationTable(db)
	}

	return db
}

// config connection database
func configConnectionDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("error load database pool, Error : %s", err)
	}

	maxIdleCon, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONNECTION"))
	if err != nil {
		return fmt.Errorf("error to get max idle connection on .env file, Error : %s", err)
	}
	MaxOpenCon, err := strconv.Atoi(os.Getenv("MAX_OPEN_CONNECTION"))
	if err != nil {
		return fmt.Errorf("error to get max open connection on .env file, Error : %s", err)
	}

	sqlDB.SetMaxIdleConns(maxIdleCon)
	sqlDB.SetMaxOpenConns(MaxOpenCon)
	sqlDB.SetConnMaxLifetime(60 * time.Minute) // set max time connection

	return nil
}

// close database connection
func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	helper.LogIfError(fmt.Errorf("failed to close connection database, Error : %s", err))

	dbSQL.Close()
}
