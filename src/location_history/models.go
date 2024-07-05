package main

import (
	"common/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)


type Location struct {
	ID		  uint       `gorm:"primaryKey;autoIncrement"`
	Username  string     `gorm:"index"`
    Longitude	  float64    
	Latitude    float64    
    Time      time.Time  `gorm:"autoCreateTime"`
}

func (loc *Location) String() string {
	return fmt.Sprintf("Coordinates: (%.8f, %.8f)",loc.Longitude,loc.Latitude)
}

// GORM hook. Executes before every save operation
func (loc *Location) BeforeSave(tx *gorm.DB) (err error) {
    loc.Longitude = utils.RoundToEightDecimals(loc.Longitude)
    loc.Latitude = utils.RoundToEightDecimals(loc.Latitude)
    return
}


func calculateDistanceByUsername(username string, startTime time.Time, endTime time.Time) (float64,error) {
	var locations []Location
	res := db.Where("Username = ? AND Time BETWEEN ? AND ?", username, startTime, endTime).Find(&locations)

	if res.Error != nil {
		return 0,res.Error
	}

	if locations == nil || len(locations) < 2{
		return 0,nil
	}

	var totalDistance float64
	prevLoc := locations[0]
	for _, currLoc := range locations[1:] {
		totalDistance += utils.CalcDistance(prevLoc.Longitude,prevLoc.Latitude,currLoc.Longitude,currLoc.Latitude)
	}

	return totalDistance,nil
}


func updateHistoryByUsername(username string, longitude float64, latitude float64) error{
	loc := Location{Username: username, Longitude: longitude, Latitude: latitude}
	db.Create(&loc)
	return nil
}



