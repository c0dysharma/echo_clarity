package helpers

import (
	"os"

	"github.com/charmbracelet/log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		log.Fatal("SQLITE_DB_PATH is not set")
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = db
	log.Info("Connected to database")

	// Migrate the schema
	// db.AutoMigrate(&Product{})
}
