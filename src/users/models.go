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
	Xcoord    float64    
	Ycoord    float64
}

func (user *User) String() string {
	return fmt.Sprintf("User[Name: %s, Coordinates: (%.8f, %.8f)]", user.Name, user.Xcoord, user.Ycoord)
}


// GORM hook, executes before each save
func (user *User) BeforeSave(tx *gorm.DB) (err error) {
    user.Xcoord = utils.RoundToEightDecimals(user.Xcoord)
    user.Ycoord = utils.RoundToEightDecimals(user.Ycoord)
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
		return res.Error
	} 

	user = User{Name: username, Xcoord: xcoord, Ycoord: ycoord}
	db.Create(&user)
	return nil
}


func getNearbyByCoordinates(xcoord, ycoord float64 , radius float64, page int) ([]User, error) {
	var users []User
	res := db.Find(&users)

	if res.Error != nil {
		return nil,res.Error
	}

	closeUsers := make([]User,0,len(users))
	for _, user := range users {
		if utils.CalcDistance(xcoord,ycoord,user.Xcoord,user.Ycoord) <= radius {
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





