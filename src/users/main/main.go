package main

import (
	"common/database"
	"common/utils"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	REST_HOST    string // Host for the REST server
	REST_PORT    string // Port for the REST server
	GRPC_HOST    string // Host for the gRPC server
	GRPC_PORT    string // Port for the gRPC server
	DATABASE_URL string // URL for the database connection
	LOG_URL      string // URL for the log file
	db           *gorm.DB // Global database connection
)

// init function loads environment variables and initializes global variables
func init() {
	REST_HOST = utils.LoadEnv("USERS_REST_HOST")
	REST_PORT = utils.LoadEnv("USERS_REST_PORT")
	GRPC_HOST = utils.LoadEnv("USERS_GRPC_HOST")
	GRPC_PORT = utils.LoadEnv("USERS_GRPC_PORT")
	DATABASE_URL = utils.LoadEnv("USERS_DATABASE_URL")
	LOG_URL = utils.LoadEnv("USERS_LOG_URL")
}

// registerRoutes registers the API routes with the Gin engine
func registerRoutes(engine *gin.Engine) {
	engine.POST("/update/:username", updateLocation)
	engine.GET("/nearby", findNearby)
}

// migrateModels migrates the database models using GORM
func migrateModels() {
	db.AutoMigrate(&User{})
}

// main function initializes logging, sets up the Gin engine, connects to the database,
// registers routes, starts the REST server, and waits for a termination signal
func main() {
	// Initialize logging to the specified log file
	file := utils.InitLogging(LOG_URL)
	defer file.Close()

	// Redirect Gin's default writer and error writer to the log file
	gin.DefaultWriter = file
	gin.DefaultErrorWriter = file

	// Create a new Gin engine and register the routes
	engine := gin.Default()
	registerRoutes(engine)

	// Connect to the database and migrate models
	db = database.New(DATABASE_URL)
	defer database.Close(db)
	migrateModels()

	// Start the REST server in a new goroutine
	go engine.Run(REST_HOST + ":" + REST_PORT)

	// Wait for a termination signal to gracefully shut down the server
	utils.WaitForSignal()
	log.Println("All services down")
}
