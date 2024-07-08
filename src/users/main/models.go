package main

import (
	"common/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const (
	PAGE_SIZE int = 3 // Constant to define the number of users per page for pagination
)

// User struct represents a user in the system with their ID, Name, Longitude, and Latitude
type User struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"` // User ID, primary key, auto-incremented
	Name      string  `gorm:"size:16;not null"`        // User Name, size limited to 16 characters, cannot be null
	Longitude float64 // User's longitude coordinate
	Latitude  float64 // User's latitude coordinate
}

// String method returns a string representation of the User struct
func (user *User) String() string {
	return fmt.Sprintf("User[Name: %s, Coordinates: (%.8f, %.8f)]", user.Name, user.Longitude, user.Latitude)
}

// BeforeSave GORM hook, executes before each save operation
// This method rounds the user's longitude and latitude to eight decimal places before saving
func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	user.Longitude = utils.RoundToEightDecimals(user.Longitude)
	user.Latitude = utils.RoundToEightDecimals(user.Latitude)
	return
}

// updateLocationByUsername updates the location of a user identified by their username
// If the user exists, it updates their longitude and latitude
// If the user does not exist, it creates a new user with the provided username, longitude, and latitude
func updateLocationByUsername(username string, longitude float64, latitude float64) error {
	var user User
	res := db.Where("Name = ?", username).First(&user)

	// If user exists, update their location
	if res.Error == nil {
		user.Longitude = longitude
		user.Latitude = latitude
		db.Save(&user)
		return nil
	}

	// If there is an error other than record not found, return the error
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	}

	// If user does not exist, create a new user
	user = User{Name: username, Longitude: longitude, Latitude: latitude}
	db.Create(&user)
	return nil
}

// getNearbyByCoordinates finds users within a certain radius from the given coordinates
// It returns a paginated list of users that are within the specified radius
func getNearbyByCoordinates(longitude float64, latitude float64, radius float64, page int) ([]User, error) {
	var users []User
	res := db.Find(&users)

	// If there is an error while fetching users, return the error
	if res.Error != nil {
		return nil, res.Error
	}

	// Filter users within the specified radius
	closeUsers := make([]User, 0, len(users))
	for _, user := range users {
		if utils.CalcDistance(longitude, latitude, user.Longitude, user.Latitude) <= radius {
			closeUsers = append(closeUsers, user)
		}
	}

	// Paginate the filtered users
	pagedUsers := make([]User, 0, PAGE_SIZE)
	firstOnPage := (page - 1) * PAGE_SIZE

	// If the requested page is beyond the available users, return an empty list
	if len(closeUsers) < firstOnPage {
		return pagedUsers, nil
	}

	// Populate the paginated list with the appropriate users
	for i := firstOnPage; i < len(closeUsers) && i < firstOnPage+PAGE_SIZE; i++ {
		pagedUsers = append(pagedUsers, closeUsers[i])
	}

	return pagedUsers, nil
}
