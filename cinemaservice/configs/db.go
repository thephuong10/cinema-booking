package configs

import (
	"cinemaservice/models/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func InitDB() *gorm.DB {

	dsn := "root:@tcp(127.0.0.1:3306)/stsrvdb?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to MySQL: " + err.Error())
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic("Failed to get raw DB object")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	autoMigrate(db)

	return db

}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&entities.Cinema{}, &entities.Room{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		log.Println("Migration successfully")
	}
}
