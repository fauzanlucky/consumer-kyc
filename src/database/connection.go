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
var dbMain *gorm.DB
var dbReplica *gorm.DB

// DB type hosts
type DB struct {
	Main    *gorm.DB
	Replica *gorm.DB
}

// dbInit Initialization Connection
// return connection, error
func dbInit(hostType string, dbName string) *gorm.DB {

	if hostType == "main" {
		hostType = os.Getenv("DB_POSTGRES_HOST_MAIN")
	} else if hostType == "replica" {
		hostType = os.Getenv("DB_POSTGRES_HOST_REPLICA")
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

// DBMain function
// return DBMain
func DBMain() *gorm.DB {
	if dbMain == nil {
		fmt.Println("No Active Main Connection Found")
		fmt.Println("Creating New Main Connection")
		dbMain = dbInit("main", os.Getenv("DB_POSTGRES_DATABASE"))
	}
	return dbMain
}

// DBReplica function
// return DBReplica
func DBReplica() *gorm.DB {
	if dbReplica == nil {
		fmt.Println("No Active Replica Connection Found")
		fmt.Println("Creating New Replica Connection")
		dbReplica = dbInit("replica", os.Getenv("DB_POSTGRES_DATABASE"))
	}
	return dbReplica
}
