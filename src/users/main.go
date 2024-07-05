package main

import (
	"common/database"
	"common/utils"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	REST_HOST string
	REST_PORT string
	GRPC_HOST string
	GRPC_PORT string
	DATABASE_URL string
	db *gorm.DB
)



func init(){
	REST_HOST = utils.LoadEnv("USERS_REST_HOST")
	REST_PORT = utils.LoadEnv("USERS_REST_PORT")
	GRPC_HOST = utils.LoadEnv("USERS_GRPC_HOST")
	GRPC_PORT = utils.LoadEnv("USERS_GRPC_PORT")
	DATABASE_URL = utils.LoadEnv("USERS_DATABASE_URL")
}



func registerRoutes(engine *gin.Engine){
	engine.POST("/update/:username",updateLocation)
	engine.GET("/nearby",findNearby)
}


func migrateModels(){
	db.AutoMigrate(&User{})
}


func main() {
	engine := gin.Default()
	registerRoutes(engine)

	var err error
	db,err = database.New(DATABASE_URL)
	if err != nil {
		fmt.Println("error: ",err.Error())
		os.Exit(1);
	}
	defer database.Close(db)
	migrateModels()
	
	engine.Run(REST_HOST + ":" + REST_PORT) 
}


