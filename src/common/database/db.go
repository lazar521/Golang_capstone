package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New creates and returns a new GORM database connection
// It configures a custom logger to log only errors and ignores the ErrRecordNotFound error
func New(url string) (*gorm.DB) {
	// Configure the custom logger
	newLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,        // Disable color
		},
	)

	// Open a new SQLite database connection with the custom logger
	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{
		Logger: newLogger,
	})

	// If there is an error opening the database, log it and exit
	if err != nil {
		log.Println("error: ", err.Error())
		os.Exit(1)
	}

	return db
}

// Close closes the given GORM database connection
// It retrieves the underlying SQL connection and closes it if it exists
func Close(db *gorm.DB) {
	if db != nil {
		dbSQL, err := db.DB()
		if err != nil {
			fmt.Println("Error getting underlying DB connection:", err)
			return
		}
		if err := dbSQL.Close(); err != nil {
			fmt.Println("Error closing database connection:", err)
		}
	}
}
