package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)



func calculateDistance(c *gin.Context){
	username := c.Param("username")
	fmt.Println("CalculateDistance username: ", username)

}