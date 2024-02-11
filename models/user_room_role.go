package models

import "gorm.io/gorm"

type UserRole string

const (
	O UserRole = "owner"
	M UserRole = "member"
)

type UserRoomRole struct {
	gorm.Model
	UserRoom int      `json:"user_room"`
	Role     UserRole `json:"role"`
}
