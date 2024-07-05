package main

import (
	"common/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)


const (
	PAGE_SIZE int = 3
)


type User struct {
    ID        uint       `gorm:"primaryKey;autoIncrement"`
    Name      string     `gorm:"size:16;not null"`
	Longitude float64    
	Latitude  float64
}

func (user *User) String() string {
	return fmt.Sprintf("User[Name: %s, Coordinates: (%.8f, %.8f)]", user.Name, user.Longitude, user.Latitude)
}


// GORM hook, executes before each save
func (user *User) BeforeSave(tx *gorm.DB) (err error) {
    user.Longitude = utils.RoundToEightDecimals(user.Longitude)
    user.Latitude = utils.RoundToEightDecimals(user.Latitude)
    return
}


func updateLocationByUsername(username string, longitude float64, latitude float64) error {
	var user User
	res := db.Where("Name = ?", username,).First(&user)
	
	if res.Error == nil {
		user.Longitude = longitude
		user.Latitude = latitude
		db.Save(&user)
		return nil
	}

	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	} 

	user = User{Name: username, Longitude: longitude, Latitude: latitude}
	db.Create(&user)
	return nil
}


func getNearbyByCoordinates(longitude float64, latitude float64 , radius float64, page int) ([]User, error) {
	var users []User
	res := db.Find(&users)

	if res.Error != nil {
		return nil,res.Error
	}

	closeUsers := make([]User,0,len(users))
	for _, user := range users {
		if utils.CalcDistance(longitude,latitude,user.Longitude,user.Latitude) <= radius {
			closeUsers = append(closeUsers, user)
		}
	}

	pagedUsers := make([]User,0,PAGE_SIZE)
	firstOnPage := (page-1)*PAGE_SIZE

	if len(closeUsers) < firstOnPage {
		return pagedUsers,nil
	} 

	for i:=firstOnPage; i<len(closeUsers) && i<firstOnPage+PAGE_SIZE; i+=1 {
		pagedUsers = append(pagedUsers, closeUsers[i])
	}

	return pagedUsers,nil
}





