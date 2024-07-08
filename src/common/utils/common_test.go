package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umahmood/haversine"
)


func TestLoadEnv(t *testing.T) {
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	value := LoadEnv("TEST_ENV")
	assert.Equal(t, "test_value", value)
}

func TestRoundToEightDecimals(t *testing.T) {
	value := 123.123456789
	expected := 123.12345679
	result := RoundToEightDecimals(value)
	assert.Equal(t, expected, result)
}

func TestCalcDistance(t *testing.T) {
	longitude1, latitude1 := 0.0, 0.0
	longitude2, latitude2 := 1.0, 1.0
	_, expectedKm := haversine.Distance(haversine.Coord{Lat: latitude1, Lon: longitude1}, haversine.Coord{Lat: latitude2, Lon: longitude2})

	result := CalcDistance(longitude1, latitude1, longitude2, latitude2)
	assert.InEpsilon(t, expectedKm, result, 0.0001, "Expected distances to be approximately equal")
}

func TestCheckUsername(t *testing.T) {
	err := CheckUsername("testuser")
	assert.NoError(t, err)

	err = CheckUsername("t")
	assert.Error(t, err, "Expected an error for username too short")

	err = CheckUsername("thisisaverylongusername")
	assert.Error(t, err, "Expected an error for username too long")

	err = CheckUsername("test_user")
	assert.Error(t, err, "Expected an error for invalid characters in username")
}

func TestCheckCoordinates(t *testing.T) {
	err := CheckCoordinates(0.0, 0.0)
	assert.NoError(t, err)

	err = CheckCoordinates(200.0, 0.0)
	assert.Error(t, err, "Expected an error for invalid longitude")

	err = CheckCoordinates(0.0, 100.0)
	assert.Error(t, err, "Expected an error for invalid latitude")
}