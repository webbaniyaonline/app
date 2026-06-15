package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	godotenv.Load(".env") // Load .env file
}

type Dbinstance struct {
	Db *gorm.DB // create DB instance
}

var DB Dbinstance

// Database connection data fetch from .env file
func ConnectDb() {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Kolkata", //Asia/Shanghai
		"localhost",
		5432,
		os.Getenv("DB_USER"),     // Get Data from ENV File
		os.Getenv("DB_PASSWORD"), // Get Data from ENV File
		os.Getenv("DB_NAME"),     // Get Data from ENV File
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		//Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err) // Display Connection Error
		//os.Exit(2)
	}

	//log.Println("connected")
	//db.Logger = logger.Default.LogMode(logger.Info)

	//log.Println("running migrations")
	//db.AutoMigrate(&models.Fact{})

	DB = Dbinstance{
		Db: db, // create DB instance
	}
}
