package models

import (
	"gorm.io/gorm"
)

type UserRole string

const (
	OWNER  UserRole = "o"
	MEMBER UserRole = "m"
)

type UserRoom struct {
	gorm.Model
	UserID string   `gorm:"not null" json:"-"`
	User   User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;not null" json:"user"`
	RoomID string   `gorm:"not null" json:"-"`
	Room   Room     `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE;not null" json:"room"`
	Role   UserRole `gorm:"size:1;not null" json:"role"` // if more than one role, so need a user_room_role model
}

func (u *UserRoom) BeforeCreate(tx *gorm.DB) error {
	u.Role = MEMBER
	return nil
}
