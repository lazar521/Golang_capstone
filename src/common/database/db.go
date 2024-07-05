package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"fmt"
)


func New(url string) (*gorm.DB,error){
	return gorm.Open(sqlite.Open(url), &gorm.Config{})
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