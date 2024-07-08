package utils

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"unicode"

	"github.com/umahmood/haversine"
)

const RADIANS_EARTH = 6371000

func InitLogging(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(file)
	file.Sync()

	return file
}


func WaitForSignal(){
	sigChan := make(chan os.Signal, 1)
	serverErrChan := make(chan error, 1)

	signal.Notify(sigChan, os.Interrupt, os.Kill)

	select {
		case sig := <-sigChan:
			log.Printf("Received signal: %v. Initiating graceful shutdown...\n", sig)
		case err := <-serverErrChan:
			log.Printf("Server error: %v. Initiating graceful shutdown...\n", err)
	}
}

func LoadEnv(variableName string) string{
	envVar := os.Getenv(variableName)
	if envVar == "" {
		fmt.Println(variableName + " env variable not provided. Exiting..")
		os.Exit(1)
	}

	return envVar
}


func RoundToEightDecimals(val float64) float64 {
    return math.Round(val*1e8) / 1e8
}


var CalcDistance = func (longitude1, latitude1, longitude2, latitude2 float64) float64 {
	coord1 := haversine.Coord{Lat: latitude1, Lon: longitude1} 
	coord2 := haversine.Coord{Lat: latitude2, Lon: longitude2} 
	_, km := haversine.Distance(coord1, coord2)
	return km
}

var CheckUsername = func(username string) error {
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

var CheckCoordinates = func(longitude, latitude float64) error{
	if longitude < -180 || longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	
	if latitude  < -90 || latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}

	return nil
}
