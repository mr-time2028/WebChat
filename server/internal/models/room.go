package models

import (
	"gorm.io/gorm"
)

type RoomType string

//const (
//	PrivateViewing RoomType = "p"
//	Group          RoomType = "g"
//)

type Room struct {
	gorm.Model
	ID   string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
	//Identifier string    `gorm:"size:36;not null" json:"identifier"`
	//Password   string    `gorm:"size:60;not null" json:"password"`
	//Type       RoomType  `gorm:"size:1;not null" json:"type"`
	Clients map[string]*Client `gorm:"-"`
}

//func (r *Room) BeforeCreate(tx *gorm.DB) error {
//	// set unique identifier
//	if r.Type == Group {
//		uuidVal, _ := uuid.NewRandom()
//		r.Identifier = uuidVal.String()
//	}
//
//	// set password
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//	r.Password = string(hashedPassword)
//
//	return nil
//}

func (r *Room) InsertOneRoom(room *Room) (string, error) {
	result := ModelRepo.db.GormDB.Create(room)
	if result.Error != nil {
		return "", result.Error
	}
	return room.ID, nil
}

func (r *Room) GetAllRooms() ([]*Room, error) {
	var rooms []*Room
	result := ModelRepo.db.GormDB.Find(&rooms)
	if result.Error != nil {
		return nil, result.Error
	}
	return rooms, nil
}
