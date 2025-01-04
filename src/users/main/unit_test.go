package main

import (
	"common/database"
	"common/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine // Global Gin engine

// wipeDatabase drops all tables and migrates the models
func wipeDatabase() {
	err := db.Migrator().DropTable(&User{})
	if err != nil {
		fmt.Println("failed to drop tables: ", err)
		os.Exit(1)
	}
	migrateModels()
}

// TestMain sets up the testing environment
func TestMain(m *testing.M) {
	// Initialize logging
	file := utils.InitLogging(LOG_URL)
	defer file.Close()

	// Set Gin's default writer and error writer to the log file
	gin.DefaultWriter = file
	gin.DefaultErrorWriter = file

	// Create a new Gin engine and register routes
	router = gin.Default()
	registerRoutes(router)

	// Connect to the database and migrate models
	db = database.New(DATABASE_URL)
	defer database.Close(db)

	// Run the tests
	m.Run()
}

// TestUpdateLocationByUsername tests the updateLocationByUsername function
func TestUpdateLocationByUsername(t *testing.T) {
	wipeDatabase()

	t.Run("Update existing user location", func(t *testing.T) {
		user := User{Name: "testuser", Longitude: 10.0, Latitude: 20.0}
		db.Create(&user)

		err := updateLocationByUsername("testuser", 30.0, 40.0)
		assert.NoError(t, err)

		var updatedUser User
		db.First(&updatedUser, "name = ?", "testuser")
		assert.Equal(t, 30.0, updatedUser.Longitude)
		assert.Equal(t, 40.0, updatedUser.Latitude)
	})

	t.Run("Create new user location", func(t *testing.T) {
		err := updateLocationByUsername("newuser", 50.0, 60.0)
		assert.NoError(t, err)

		var newUser User
		db.First(&newUser, "name = ?", "newuser")
		assert.Equal(t, "newuser", newUser.Name)
		assert.Equal(t, 50.0, newUser.Longitude)
		assert.Equal(t, 60.0, newUser.Latitude)
	})
}

// TestGetNearbyByCoordinates tests the getNearbyByCoordinates function
func TestGetNearbyByCoordinates(t *testing.T) {
	wipeDatabase()

	users := []User{
		{Name: "user1", Longitude: 10.0, Latitude: 10.0},
		{Name: "user2", Longitude: 20.0, Latitude: 20.0},
		{Name: "user3", Longitude: 30.0, Latitude: 30.0},
		{Name: "user4", Longitude: 40.0, Latitude: 40.0},
	}
	db.Create(&users)

	t.Run("Find users within radius", func(t *testing.T) {
		nearbyUsers, err := getNearbyByCoordinates(15.0, 15.0, 2000.0, 1)
		assert.NoError(t, err)
		assert.Len(t, nearbyUsers, 2)
		assert.Equal(t, nearbyUsers[0], users[0])
		assert.Equal(t, nearbyUsers[1], users[1])
	})

	t.Run("No users within radius", func(t *testing.T) {
		nearbyUsers, err := getNearbyByCoordinates(0.0, 0.0, 5.0, 1)
		assert.NoError(t, err)
		assert.Len(t, nearbyUsers, 0)
	})

	t.Run("Pagination test", func(t *testing.T) {
		nearbyUsers, err := getNearbyByCoordinates(15.0, 15.0, 100000.0, 2)
		assert.NoError(t, err)
		assert.Len(t, nearbyUsers, 1)
	})
}

// TestUpdateLocation tests the updateLocation endpoint
func TestUpdateLocation(t *testing.T) {
	wipeDatabase()
	
	// Mock the notifyLocationHistoryService function
	notifyLocationHistoryService = func(username string, longitude, latitude float64) error {
		return nil
	}

	users := []User{
		{Name: "testuser", Longitude: 0.0, Latitude: 10.0},
	}
	db.Create(&users)

	t.Run("Valid Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/update/testuser", strings.NewReader(`{"longitude": 10.0, "latitude": 20.0}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"Username": "testuser", "Longitude": 10.0, "Latitude": 20.0}`, w.Body.String())
	})

	t.Run("Invalid Username", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/update/", strings.NewReader(`{"longitude": 10.0, "latitude": 20.0}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Invalid Coordinates", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/update/testuser", strings.NewReader(`{"longitude": 200.0, "latitude": 100.0}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "longitude must be between -180 and 180"}`, w.Body.String())
	})
}
