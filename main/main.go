package main

import (
	"github.com/gin-gonic/gin"
)


func registerRoutes(engine *gin.Engine){
	engine.POST("/update/:username")
	engine.GET("/nearby")
	engine.GET("/distance/:username")
}



func main() {
	engine := gin.Default()
	
	registerRoutes(engine)

	engine.Run() // listen and serve on 0.0.0.0:8080
}


