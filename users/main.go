package main

import (
	"github.com/gin-gonic/gin"
)




func registerRoutes(engine *gin.Engine){
	engine.POST("/update/:username",updateLocation)
	engine.GET("/nearby",findNearby)
}



func main() {
	engine := gin.Default()
	registerRoutes(engine)
	engine.Run("localhost:8001") 
}


