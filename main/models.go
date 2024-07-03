package main

import (
	"fmt"
	"os"
	"time"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math"
)

var db *gorm.DB

func init(){
	var err error
	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error occured: ", err)
		os.Exit(1)
	}

	db.AutoMigrate(&User{}, &Location{})
}



func roundToEightDecimals(val float64) float64 {
    return math.Round(val*1e8) / 1e8
}

func (loc *Location) BeforeSave(tx *gorm.DB) (err error) {
    loc.Xcoord = roundToEightDecimals(loc.Xcoord)
    loc.Ycoord = roundToEightDecimals(loc.Ycoord)
    return
}





type User struct {
    ID        uint       `gorm:"primaryKey;autoIncrement"`
    Name      string     `gorm:"size:100;not null"`
	Locations []Location `gorm:"foreignKey:UserRefer"`
}



type Location struct {
	UserRefer uint       `gorm:"index"`
    Xcoord	  float64    
	Ycoord    float64    
    Time      time.Time  `gorm:"autoCreateTime"`
}