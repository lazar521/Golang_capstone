package main

import (
	"fmt"
	"math"
	"os"
	"time"

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

	databaseURL = databaseURL + "/location_history.db"

	var err error
	db, err = gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		fmt.Println("Error occured: ", err)
		os.Exit(1)
	}

	db.AutoMigrate( &Location{})
}



func roundToEightDecimals(val float64) float64 {
    return math.Round(val*1e8) / 1e8
}

func (loc *Location) BeforeSave(tx *gorm.DB) (err error) {
    loc.Xcoord = roundToEightDecimals(loc.Xcoord)
    loc.Ycoord = roundToEightDecimals(loc.Ycoord)
    return
}




type Location struct {
	UserRefer uint       `gorm:"index"`
    Xcoord	  float64    
	Ycoord    float64    
    Time      time.Time  `gorm:"autoCreateTime"`
}