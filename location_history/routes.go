package main

import (
	"net/http"

	"common/utils"

	"github.com/gin-gonic/gin"
)



func getTraveledDistance(c *gin.Context){
	username := c.Param("username")
	
	if err := utils.CheckUsername(username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err})
		return
	}

	distance,err := calculateDistanceByUsername(username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Traveled distance":distance})
}


