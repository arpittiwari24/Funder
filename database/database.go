package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"main.go/models"
)

type DBInstance struct {
	DB *gorm.DB
}

var DB DBInstance

func ConnectDb()  {

	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Kolkata",
                   os.Getenv("DB_USER"),
                   os.Getenv("DB_PASSWORD"),
                   os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	fmt.Println(("This line reached upar wali"))
	if err != nil {
		log.Fatal("Failed to connect\n",err)
		os.Exit(2)
	}
	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	db.AutoMigrate(&models.Product{})
	fmt.Println(("This line reached"))
	DB = DBInstance{
		DB: db,
	}
}