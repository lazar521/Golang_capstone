package main

import (
	"common/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	LAYOUT string = time.RFC3339 // Time layout for parsing and formatting
)

// getTraveledDistance handles the HTTP GET request to calculate the distance traveled by a user
func getTraveledDistance(c *gin.Context) {
	username := c.Param("username")

	// Check if the username is valid
	if err := utils.CheckUsername(username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Struct to bind query parameters
	data := struct {
		StartTimeStr string `form:"start"`
		EndTimeStr   string `form:"end"`
	}{}

	// Bind the query parameters to the struct
	if err := c.ShouldBindQuery(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startExists := (data.StartTimeStr != "")
	endExists := (data.EndTimeStr != "")
	// Ensure both start and end times are provided or none at all
	if (startExists && !endExists) || (!startExists && endExists) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide either both lower and upper time bound or none"})
		return
	}

	// Replace spaces with plus signs in the time strings
	data.StartTimeStr = strings.ReplaceAll(data.StartTimeStr, " ", "+")
	data.EndTimeStr = strings.ReplaceAll(data.EndTimeStr, " ", "+")

	var startTime time.Time
	var endTime time.Time

	if startExists {
		var err error
		startTime, err = time.Parse(LAYOUT, data.StartTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "lower time bound of unknown format"})
			return
		}

		endTime, err = time.Parse(LAYOUT, data.EndTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "upper time bound of unknown format"})
			return
		}

		// Ensure the start time is before the end time
		if startTime.After(endTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end time is set before start time"})
			return
		}
	} else {
		// Default to the last 24 hours if no time bounds are provided
		endTime = time.Now()
		startTime = endTime.AddDate(0, 0, -1)
	}

	// Calculate the total distance traveled by the user
	distance, err := calculateDistanceByUsername(username, startTime, endTime)
	if err != nil {
		log.Println("error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not calculate distance"})
		return
	}

	// Return the calculated distance
	c.JSON(http.StatusOK, gin.H{"Traveled distance": distance})
}
