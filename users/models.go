package main

import (
	"fmt"
	"math"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


var db *gorm.DB

func init(){
	databaseURL := os.Getenv("DATA_FOLDER")
	if databaseURL == "" {
		fmt.Println("DATA_FOLDER env variable not provided. Exiting..")
		os.Exit(1)
	}

	databaseURL = databaseURL + "/users.db"

	var err error
	db, err = gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		fmt.Println("Error occured: ", err)
		os.Exit(1)
	}

	db.AutoMigrate(&User{})
}


func roundToEightDecimals(val float64) float64 {
    return math.Round(val*1e8) / 1e8
}

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
    user.Xcoord = roundToEightDecimals(user.Xcoord)
    user.Ycoord = roundToEightDecimals(user.Ycoord)
    return
}



type User struct {
    ID        uint       `gorm:"primaryKey;autoIncrement"`
    Name      string     `gorm:"size:16;not null"`
	Xcoord    float64    
	Ycoord    float64
}


