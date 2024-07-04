package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)




func updateLocation(c *gin.Context){
	data := struct {
		Xcoord float64 `form:"xcoord" binding:"required"`
		Ycoord float64 `form:"ycoord" binding:"required"`
	}{}
	
	username := c.Param("username")
	
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := checkUsername(username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := checkCoordinates(data.Xcoord,data.Ycoord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := updateLocationByUsername(username,data.Xcoord,data.Ycoord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"you username" : username, "Xcoord":data.Xcoord,"Ycoord":data.Ycoord})
}



func findNearby(c *gin.Context){
	data := struct{
		Xcoord float64 `form:"xcoord" binding:"required"`
		Ycoord float64 `form:"ycoord" binding:"required"`
		Radius float64 `form:"radius" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := checkCoordinates(data.Xcoord,data.Ycoord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := getNearbyByCoordinates(data.Xcoord,data.Ycoord,data.Radius);
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Closeby" : users,
    })
}
