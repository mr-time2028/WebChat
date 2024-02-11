package models

import "gorm.io/gorm"

type UserRoomPermission struct {
	gorm.Model
	UserRoom   int `json:"user_room"`
	Permission int `json:"permission"`
}
