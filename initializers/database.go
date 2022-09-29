package initializers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB")
	config := &gorm.Config{}
	config.Logger = logger.Default.LogMode(logger.Info)

	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

}
