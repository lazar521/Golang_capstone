package main

import (
	"common/database"
	"common/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine // Global Gin engine

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

	// Drop existing tables and migrate models
	err := db.Migrator().DropTable(&Location{})
	if err != nil {
		fmt.Println("failed to drop tables: ", err)
		os.Exit(1)
	}
	migrateModels()

	// Run the tests
	m.Run()
}

// TestCalculateDistanceByUsername tests the calculateDistanceByUsername function
func TestCalculateDistanceByUsername(t *testing.T) {
	expected := 30.507941089707185

	locations := []Location{
		{Username: "testuser", Longitude: 10.0, Latitude: 20.0, Time: time.Now().Add(-10 * time.Minute)},
		{Username: "testuser", Longitude: 10.1, Latitude: 20.1, Time: time.Now().Add(-5 * time.Minute)},
		{Username: "testuser", Longitude: 10.2, Latitude: 20.2, Time: time.Now()},
	}

	for _, loc := range locations {
		db.Create(&loc)
	}

	startTime := time.Now().Add(-15 * time.Minute)
	endTime := time.Now()

	// Calculate the distance traveled by the user
	distance, err := calculateDistanceByUsername("testuser", startTime, endTime)

	assert.NoError(t, err)
	assert.Equal(t, expected, distance)
}

// TestUpdateHistoryByUsername tests the updateHistoryByUsername function
func TestUpdateHistoryByUsername(t *testing.T) {
	// Update the location history for the user
	err := updateHistoryByUsername("testuser", 10.0, 20.0)
	assert.NoError(t, err)

	// Retrieve the location from the database and check its values
	var location Location
	db.Where("Username = ?", "testuser").First(&location)

	assert.Equal(t, "testuser", location.Username)
	assert.Equal(t, 10.0, location.Longitude)
	assert.Equal(t, 20.0, location.Latitude)
}

// TestGetTraveledDistance tests the getTraveledDistance endpoint
func TestGetTraveledDistance(t *testing.T) {
	t.Run("Valid Request with Time Bounds", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/distance/testuser?start=2022-07-08T00:00:00Z&end=2025-07-09T00:00:00Z", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"Traveled distance": 61.01587896204506}`, w.Body.String())
	})

	t.Run("Valid Request without Time Bounds", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/distance/testuser", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"Traveled distance": 61.01587896204506}`, w.Body.String())
	})

	t.Run("Invalid Username", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/distance/invaliduser", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"Traveled distance": 0}`, w.Body.String())
	})

	t.Run("Invalid Time Bounds", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/distance/testuser?start=invalid&end=2023-07-09T00:00:00Z", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"lower time bound of unknown format"}`, w.Body.String())
	})

	t.Run("Start Time without End Time", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/distance/testuser?start=2023-07-08T00:00:00Z", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"provide either both lower and upper time bound or none"}`, w.Body.String())
	})
}
