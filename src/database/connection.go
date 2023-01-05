package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB Global Connection
var dbMaster *gorm.DB
var dbSlave *gorm.DB

// DB type hosts
type DB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

// dbInit Initialization Connection
// return connection, error
func dbInit(hostType string, dbName string) *gorm.DB {

	if hostType == "master" {
		hostType = os.Getenv("DB_POSTGRES_HOST_MASTER")
	} else if hostType == "slave" {
		hostType = os.Getenv("DB_POSTGRES_HOST_SLAVE")
	}

	postgresCon := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		hostType,
		os.Getenv("DB_POSTGRES_PORT"),
		os.Getenv("DB_POSTGRES_USERNAME"),
		dbName,
		os.Getenv("DB_POSTGRES_PASSWORD"),
	)

	DB, err := gorm.Open(postgres.Open(postgresCon), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Println(err.Error())
		panic(fmt.Sprintf("Failed to connect to database %s", postgresCon))
	}

	fmt.Printf("Successfully connected to database %s\n", postgresCon)

	connConfDB, err := DB.DB()

	if err != nil {
		panic(fmt.Sprintf("Failed to set connection configuration %s", err))
	}

	connConfDB.SetConnMaxLifetime(5 * time.Minute)
	connConfDB.SetMaxIdleConns(20)
	connConfDB.SetMaxOpenConns(200)
	connConfDB.SetMaxOpenConns(20)

	fmt.Println("Connection is created")

	return DB
}

// DBMaster function
// return DBMaster
func DBMaster() *gorm.DB {
	if dbMaster == nil {
		fmt.Println("No Active Master Connection Found")
		fmt.Println("Creating New Master Connection")
		dbMaster = dbInit("master", os.Getenv("DB_POSTGRES_DATABASE"))
	}
	return dbMaster
}

// DBSlave function
// return DBSlave
func DBSlave() *gorm.DB {
	if dbSlave == nil {
		fmt.Println("No Active Slave Connection Found")
		fmt.Println("Creating New Slave Connection")
		dbSlave = dbInit("slave", os.Getenv("DB_POSTGRES_DATABASE"))
	}
	return dbSlave
}
