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
    Xcoord	  float64    
	Ycoord    float64    
    Time      time.Time  `gorm:"autoCreateTime"`
}

func (loc *Location) String() string {
	return fmt.Sprintf("Coordinates: (%.8f, %.8f)",loc.Xcoord,loc.Ycoord)
}

// GORM hook. Executes before every save operation
func (loc *Location) BeforeSave(tx *gorm.DB) (err error) {
    loc.Xcoord = utils.RoundToEightDecimals(loc.Xcoord)
    loc.Ycoord = utils.RoundToEightDecimals(loc.Ycoord)
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
		totalDistance += utils.CalcDistance(prevLoc.Xcoord,prevLoc.Ycoord,currLoc.Xcoord,currLoc.Ycoord)
	}

	return totalDistance,nil
}


func updateHistoryByUsername(username string, xcoord float64, ycoord float64) error{
	loc := Location{Username: username, Xcoord: xcoord, Ycoord: ycoord}
	db.Create(&loc)
	return nil
}



