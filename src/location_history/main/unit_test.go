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


var router *gin.Engine

// setup
func TestMain(m *testing.M) {
	file := utils.InitLogging(LOG_URL)
	defer 	file.Close()

	gin.DefaultWriter = file
	gin.DefaultErrorWriter = file

	router = gin.Default()
	registerRoutes(router)

	db = database.New(DATABASE_URL)
	defer database.Close(db)

	err := db.Migrator().DropTable(&Location{})
    if err != nil {
        fmt.Println("failed to drop tables: ", err)
		os.Exit(1)
	}

	migrateModels()

    m.Run()
}



func TestCalculateDistanceByUsername(t *testing.T) {
	expected := 45.76419124040538
	
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

	distance, err := calculateDistanceByUsername("testuser", startTime, endTime)
	
	assert.NoError(t, err)
	assert.Equal(t, distance, expected)
}


func TestUpdateHistoryByUsername(t *testing.T) {
	err := updateHistoryByUsername("testuser", 10.0, 20.0)
	assert.NoError(t, err)

	var location Location
	db.Where("Username = ?", "testuser").First(&location)

	assert.Equal(t, "testuser", location.Username)
	assert.Equal(t, 10.0, location.Longitude)
	assert.Equal(t, 20.0, location.Latitude)
}




func TestGetTraveledDistance(t *testing.T) {
    t.Run("Valid Request with Time Bounds", func(t *testing.T) {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/distance/testuser?start=2022-07-08T00:00:00Z&end=2025-07-09T00:00:00Z", nil)
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
        assert.JSONEq(t, `{"Traveled distance": 45.76419124040538}`, w.Body.String())
    })

    t.Run("Valid Request without Time Bounds", func(t *testing.T) {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/distance/testuser", nil)
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
        assert.JSONEq(t, `{"Traveled distance":45.76419124040538}`, w.Body.String())
    })

    t.Run("Invalid Username", func(t *testing.T) {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/distance/invaliduser", nil)
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
        assert.JSONEq(t, `{"Traveled distance":0}`, w.Body.String())
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