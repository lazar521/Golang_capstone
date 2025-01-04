package main

import (
	"common/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// updateLocation handles the HTTP POST request to update a user's location.
// It validates the request parameters, updates the user's location in the database,
// and notifies the location history service.
func updateLocation(c *gin.Context) {
	// Struct to bind JSON request data
	data := struct {
		Longitude float64 `form:"longitude" binding:"required"`
		Latitude  float64 `form:"latitude" binding:"required"`
	}{}

	// Get the username from the URL parameter
	username := c.Param("username")

	// Bind the JSON request data to the struct
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the username is valid
	if err := utils.CheckUsername(username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the coordinates are valid
	if err := utils.CheckCoordinates(data.Longitude, data.Latitude); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user's location in the database
	if err := updateLocationByUsername(username, data.Longitude, data.Latitude); err != nil {
		log.Println("Error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user location and history"})
		return
	}

	// Notify the location history service of the updated location
	if err := notifyLocationHistoryService(username, data.Longitude, data.Latitude); err != nil {
		log.Println("Error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update location history"})
		return
	}

	// Return a successful response
	c.JSON(http.StatusOK, gin.H{"Username": username, "Longitude": data.Longitude, "Latitude": data.Latitude})
}

// findNearby handles the HTTP GET request to find nearby users.
// It validates the request parameters, retrieves the nearby users from the database,
// and returns the results.
func findNearby(c *gin.Context) {
	// Struct to bind query parameters
	data := struct {
		Longitude float64 `form:"longitude" binding:"required"`
		Latitude  float64 `form:"latitude" binding:"required"`
		Radius    float64 `form:"radius" binding:"required"`
		Page      int     `form:"page" binding:"required"`
	}{}

	// Bind the query parameters to the struct
	if err := c.ShouldBindQuery(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the page number is valid
	if data.Page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page number must be greater than zero"})
		return
	}

	// Check if the coordinates are valid
	if err := utils.CheckCoordinates(data.Longitude, data.Latitude); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the nearby users from the database
	users, err := getNearbyByCoordinates(data.Longitude, data.Latitude, data.Radius, data.Page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the list of nearby users
	c.JSON(http.StatusOK, gin.H{
		"Closeby": users,
	})
}
