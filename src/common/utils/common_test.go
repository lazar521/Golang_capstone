package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umahmood/haversine"
)

// TestLoadEnv tests the LoadEnv function
// It sets an environment variable, verifies its value using LoadEnv, and then unsets it
func TestLoadEnv(t *testing.T) {
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	value := LoadEnv("TEST_ENV")
	assert.Equal(t, "test_value", value)
}

// TestRoundToEightDecimals tests the RoundToEightDecimals function
// It verifies that a value is correctly rounded to eight decimal places
func TestRoundToEightDecimals(t *testing.T) {
	value := 123.12345678912
	expected := 123.12345679
	result := RoundToEightDecimals(value)
	assert.Equal(t, expected, result)
}

// TestCalcDistance tests the CalcDistance function
// It verifies that the calculated distance between two coordinates matches the expected value
func TestCalcDistance(t *testing.T) {
	longitude1, latitude1 := 0.0, 0.0
	longitude2, latitude2 := 1.0, 1.0
	_, expectedKm := haversine.Distance(haversine.Coord{Lat: latitude1, Lon: longitude1}, haversine.Coord{Lat: latitude2, Lon: longitude2})

	result := CalcDistance(longitude1, latitude1, longitude2, latitude2)
	assert.InEpsilon(t, expectedKm, result, 0.0001, "Expected distances to be approximately equal")
}

// TestCheckUsername tests the CheckUsername function
// It verifies that the function correctly validates usernames based on length and character criteria
func TestCheckUsername(t *testing.T) {
	// Test a valid username
	err := CheckUsername("testuser")
	assert.NoError(t, err)

	// Test a username that is too short
	err = CheckUsername("t")
	assert.Error(t, err, "Expected an error for username too short")

	// Test a username that is too long
	err = CheckUsername("thisisaverylongusername")
	assert.Error(t, err, "Expected an error for username too long")

	// Test a username with invalid characters
	err = CheckUsername("test_user")
	assert.Error(t, err, "Expected an error for invalid characters in username")
}

// TestCheckCoordinates tests the CheckCoordinates function
// It verifies that the function correctly validates geographical coordinates
func TestCheckCoordinates(t *testing.T) {
	// Test valid coordinates
	err := CheckCoordinates(0.0, 0.0)
	assert.NoError(t, err)

	// Test invalid longitude
	err = CheckCoordinates(200.0, 0.0)
	assert.Error(t, err, "Expected an error for invalid longitude")

	// Test invalid latitude
	err = CheckCoordinates(0.0, 100.0)
	assert.Error(t, err, "Expected an error for invalid latitude")
}
