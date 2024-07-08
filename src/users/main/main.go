package main

import (
	"common/database"
	"common/utils"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	REST_HOST string
	REST_PORT string
	GRPC_HOST string
	GRPC_PORT string
	DATABASE_URL string
	LOG_URL string
	db *gorm.DB

)


func init(){
	REST_HOST = utils.LoadEnv("USERS_REST_HOST")
	REST_PORT = utils.LoadEnv("USERS_REST_PORT")
	GRPC_HOST = utils.LoadEnv("USERS_GRPC_HOST")
	GRPC_PORT = utils.LoadEnv("USERS_GRPC_PORT")
	DATABASE_URL = utils.LoadEnv("USERS_DATABASE_URL")
	LOG_URL = utils.LoadEnv("USERS_LOG_URL")
}


func registerRoutes(engine *gin.Engine){
	engine.POST("/update/:username",updateLocation)
	engine.GET("/nearby",findNearby)
}


func migrateModels(){
	db.AutoMigrate(&User{})
}




func main() {
	file := utils.InitLogging(LOG_URL)
	defer 	file.Close()
	
	gin.DefaultWriter = file
	gin.DefaultErrorWriter = file

	engine := gin.Default()
	registerRoutes(engine)

	db = database.New(DATABASE_URL)
	defer database.Close(db)
	migrateModels()
	
	go engine.Run(REST_HOST + ":" + REST_PORT) 

	utils.WaitForSignal()
	log.Println("All services down")
}


