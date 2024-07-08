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


func New(url string) (*gorm.DB){
    newLogger := logger.New(
        log.Default(),
        logger.Config{
            SlowThreshold:             time.Second, // Slow SQL threshold
            LogLevel:                  logger.Error, // Log level
            IgnoreRecordNotFoundError: true, // Ignore ErrRecordNotFound error for logger
            Colorful:                  false, // Disable color
        },
    )

	db,err := gorm.Open(sqlite.Open(url), &gorm.Config{
        Logger: newLogger,
    })

    if err != nil {
		log.Println("error: ",err.Error())
		os.Exit(1);
	}

    return db
}


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