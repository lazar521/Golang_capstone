package main

import (
	"common/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)



func updateLocation(c *gin.Context){
	data := struct {
		Longitude float64 `form:"longitude" binding:"required"`
		Latitude float64 `form:"latitude" binding:"required"`
	}{}
	
	username := c.Param("username")
	
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.CheckUsername(username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := utils.CheckCoordinates(data.Longitude,data.Latitude); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	if err := updateLocationByUsername(username,data.Longitude,data.Latitude); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := notifyLocationHistoryService(username,data.Longitude,data.Latitude); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Username" : username, "Longitude":data.Longitude,"Latitude":data.Latitude})
}



func findNearby(c *gin.Context){
	data := struct{
		Longitude float64 `form:"longitude" binding:"required"`
		Latitude float64  `form:"latitude" binding:"required"`
		Radius float64    `form:"radius" binding:"required"`
		Page   int        `form:"page"   binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page number must be greater than zero"})
		return
	}

	if err := utils.CheckCoordinates(data.Longitude,data.Latitude); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := getNearbyByCoordinates(data.Longitude,data.Latitude,data.Radius,data.Page);
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Closeby" : users,
    })
}
