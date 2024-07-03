package main

import (
	"github.com/gin-gonic/gin"
)




func registerRoutes(engine *gin.Engine){
	engine.GET("/distance/:username",calculateDistance)
}



func main() {
	engine := gin.Default()
	registerRoutes(engine)
	engine.Run("localhost:8000") 
}


