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
	LAYOUT string = time.RFC3339
)


func getTraveledDistance(c *gin.Context){
	username := c.Param("username")

	if err := utils.CheckUsername(username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err})
		return
	}

	data := struct{
		StartTimeStr string `form:"start"`
		EndTimeStr   string `form:"end"`
	}{}

	if err := c.ShouldBindQuery(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startExists := (data.StartTimeStr != "")
	endExists := (data.EndTimeStr != "")
	if (startExists && !endExists) || (!startExists && endExists){
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide either both lower and upper time bound or none"})
		return
	}

	data.StartTimeStr = strings.ReplaceAll(data.StartTimeStr, " ", "+")
	data.EndTimeStr = strings.ReplaceAll(data.EndTimeStr, " ", "+")

	var startTime time.Time
	var endTime time.Time

	if startExists {
		var err error
		startTime, err = time.Parse(LAYOUT, data.StartTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "lower time bound of unknown format" + startTime.String()})
			return
		}
	
		endTime, err = time.Parse(LAYOUT, data.EndTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "upper time bound of unknown format" + endTime.String()})
			return
		}
	} else {
		endTime = time.Now()
		startTime = endTime.AddDate(0,0,-1)
	}

	distance,err := calculateDistanceByUsername(username,startTime,endTime)
	if err != nil {
		log.Println("error: ",err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error":"could not calculate distance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Traveled distance":distance})
}


