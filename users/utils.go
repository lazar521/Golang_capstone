package main

import (
	"errors"
	"math"
	"unicode"

	"github.com/umahmood/haversine"
)


const RADIANS_EARTH = 6371000

func roundToEightDecimals(val float64) float64 {
    return math.Round(val*1e8) / 1e8
}


func calcDistance(xcoord1, ycoord1, xcoord2, ycoord2 float64) float64 {
	coord1 := haversine.Coord{Lat: ycoord1, Lon: xcoord1} 
	coord2 := haversine.Coord{Lat: ycoord2, Lon: xcoord2} 
	_, km := haversine.Distance(coord1, coord2)
	return km
}

func checkUsername(username string) error {
	if len(username) < 4 || len(username) > 16 {
		return errors.New("username must be between 4 and 16 characters long")
	}

	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return errors.New("username can only contain letters (a-z, A-Z) and numbers (0-9)")
		}
	}

	return nil
}

func checkCoordinates(xcoord, ycoord float64) error{
	if xcoord < -180 || xcoord > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	
	if ycoord  < -90 || ycoord > 90 {
		return errors.New("latitude must be between -90 and 90")
	}

	return nil
}