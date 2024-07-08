package main

import (
	"common/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Location represents a geographical location with a username, coordinates, and timestamp
type Location struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"` // Primary key, auto-incremented
	Username  string    `gorm:"index"`                   // Indexed username
	Longitude float64   // Longitude coordinate
	Latitude  float64   // Latitude coordinate
	Time      time.Time `gorm:"autoCreateTime"`          // Timestamp, auto-created on record insertion
}

// String method returns a string representation of the Location struct
func (loc *Location) String() string {
	return fmt.Sprintf("Coordinates: (%.8f, %.8f)", loc.Longitude, loc.Latitude)
}

// BeforeSave GORM hook, executes before every save operation
// This method rounds the longitude and latitude to eight decimal places before saving
func (loc *Location) BeforeSave(tx *gorm.DB) (err error) {
	loc.Longitude = utils.RoundToEightDecimals(loc.Longitude)
	loc.Latitude = utils.RoundToEightDecimals(loc.Latitude)
	return
}

// calculateDistanceByUsername calculates the total distance traveled by a user between two timestamps
// It retrieves the user's locations from the database and sums up the distances between consecutive points
func calculateDistanceByUsername(username string, startTime time.Time, endTime time.Time) (float64, error) {
	var locations []Location
	res := db.Where("Username = ? AND Time BETWEEN ? AND ?", username, startTime, endTime).Find(&locations)

	if res.Error != nil {
		return 0, res.Error
	}

	if locations == nil || len(locations) < 2 {
		return 0, nil
	}

	var totalDistance float64
	prevLoc := locations[0]
	for _, currLoc := range locations[1:] {
		totalDistance += utils.CalcDistance(prevLoc.Longitude, prevLoc.Latitude, currLoc.Longitude, currLoc.Latitude)
		prevLoc = currLoc
	}

	return totalDistance, nil
}

// updateHistoryByUsername updates the location history for a given username
// It creates a new Location record with the provided username, longitude, and latitude
func updateHistoryByUsername(username string, longitude float64, latitude float64) error {
	loc := Location{Username: username, Longitude: longitude, Latitude: latitude}
	db.Create(&loc)
	return nil
}
