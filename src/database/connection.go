package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/forkyid/consumer-kyc-update/src/constant"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB Global Connection
var dbCMSMain *gorm.DB
var dbCMSReplica *gorm.DB
var dbGiftshopMain *gorm.DB
var dbGiftshopReplica *gorm.DB

// DB type hosts
type DB struct {
	CMSMain         *gorm.DB
	CMSReplica      *gorm.DB
	GiftshopMain    *gorm.DB
	GiftshopReplica *gorm.DB
}

// dbInit Initialization Connection
// return connection, error
func dbInit(hostType string, dbName string) *gorm.DB {
	var port, username, password string

	switch hostType {
	case constant.DBHostTypeCMSMain:
		hostType = os.Getenv("DB_POSTGRES_HOST_CMS_MAIN")
		port = os.Getenv("DB_POSTGRES_CMS_PORT")
		username = os.Getenv("DB_POSTGRES_CMS_USERNAME")
		password = os.Getenv("DB_POSTGRES_CMS_PASSWORD")
	case constant.DBHostTypeCMSReplica:
		hostType = os.Getenv("DB_POSTGRES_HOST_CMS_REPLICA")
		port = os.Getenv("DB_POSTGRES_CMS_PORT")
		username = os.Getenv("DB_POSTGRES_CMS_USERNAME")
		password = os.Getenv("DB_POSTGRES_CMS_PASSWORD")
	case constant.DBHostTypeGSMain:
		hostType = os.Getenv("DB_POSTGRES_HOST_GIFTSHOP_MAIN")
		port = os.Getenv("DB_POSTGRES_GIFTSHOP_PORT")
		username = os.Getenv("DB_POSTGRES_GIFTSHOP_USERNAME")
		password = os.Getenv("DB_POSTGRES_GIFTSHOP_PASSWORD")
	case constant.DBHostTypeGSReplica:
		hostType = os.Getenv("DB_POSTGRES_HOST_GIFTSHOP_REPLICA")
		port = os.Getenv("DB_POSTGRES_GIFTSHOP_PORT")
		username = os.Getenv("DB_POSTGRES_GIFTSHOP_USERNAME")
		password = os.Getenv("DB_POSTGRES_GIFTSHOP_PASSWORD")
	}

	postgresCon := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		hostType,
		port,
		username,
		dbName,
		password,
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

	fmt.Println("Connection is created")

	return DB
}

func DBCMSMain() *gorm.DB {
	if dbCMSMain == nil {
		fmt.Println("No Active Main Connection Found")
		fmt.Println("Creating New Main Connection")
		dbCMSMain = dbInit(constant.DBHostTypeCMSMain, os.Getenv("DB_POSTGRES_DATABASE_CMS"))
	}
	return dbCMSMain
}

func DBCMSReplica() *gorm.DB {
	if dbCMSReplica == nil {
		fmt.Println("No Active Replica Connection Found")
		fmt.Println("Creating New Replica Connection")
		dbCMSReplica = dbInit(constant.DBHostTypeCMSReplica, os.Getenv("DB_POSTGRES_DATABASE_CMS"))
	}
	return dbCMSReplica
}

func DBGiftshopMain() *gorm.DB {
	if dbGiftshopMain == nil {
		fmt.Println("No Active Main Connection Found")
		fmt.Println("Creating New Main Connection")
		dbGiftshopMain = dbInit(constant.DBHostTypeGSMain, os.Getenv("DB_POSTGRES_DATABASE_GS"))
	}
	return dbGiftshopMain
}

func DBGiftshopReplica() *gorm.DB {
	if dbGiftshopReplica == nil {
		fmt.Println("No Active Replica Connection Found")
		fmt.Println("Creating New Replica Connection")
		dbGiftshopReplica = dbInit(constant.DBHostTypeGSReplica, os.Getenv("DB_POSTGRES_DATABASE_GS"))
	}
	return dbGiftshopReplica
}
