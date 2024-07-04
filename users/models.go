package main

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


var db *gorm.DB

func init(){
	databaseURL := os.Getenv("DATA_FOLDER")
	if databaseURL == "" {
		fmt.Println("DATA_FOLDER env variable not provided. Exiting..")
		os.Exit(1)
	}

	databaseURL = databaseURL + "/users.db"

	var err error
	db, err = gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		fmt.Println("Error occured: ", err)
		os.Exit(1)
	}

	db.AutoMigrate(&User{})
}


// GORM hook, executes before each save
func (user *User) BeforeSave(tx *gorm.DB) (err error) {
    user.Xcoord = roundToEightDecimals(user.Xcoord)
    user.Ycoord = roundToEightDecimals(user.Ycoord)
    return
}


func updateLocationByUsername(username string, xcoord float64, ycoord float64) error {
	var user User
	res := db.Where("Name = ?", username,).First(&user)
	
	if res.Error == nil {
		user.Xcoord = xcoord
		user.Ycoord = ycoord
		db.Save(&user)
		return nil
	}

	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("Error occurred while querying the database:", res.Error)
		return res.Error
	} 

	user = User{Name: username, Xcoord: xcoord, Ycoord: ycoord}
	db.Create(&user)
	return nil
}


func getNearbyByCoordinates(xcoord, ycoord, radius float64) ([]User, error) {
	var users []User
	res := db.Find(&users)

	if res.Error != nil {
		return nil,res.Error
	}

	closeUsers := make([]User,0,len(users))

	for _, user := range users {
		if calcDistance(xcoord,ycoord,user.Xcoord,user.Ycoord) <= radius {
			closeUsers = append(closeUsers, user)
		}
	}

	return closeUsers,nil
}


func (user *User) String() string {
	return fmt.Sprintf("User[Name: %s, Coordinates: (%.8f, %.8f)]", user.Name, user.Xcoord, user.Ycoord)
}


type User struct {
    ID        uint       `gorm:"primaryKey;autoIncrement"`
    Name      string     `gorm:"size:16;not null"`
	Xcoord    float64    
	Ycoord    float64
}


