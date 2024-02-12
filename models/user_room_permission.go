package models

import "gorm.io/gorm"

type UserRoomPermission struct {
	gorm.Model
	UserRoomID   int        `gorm:"not null" json:"-"`
	UserRoom     UserRoom   `gorm:"foreignKey:UserRoomID;constraint:OnDelete:CASCADE;not null" json:"user_room"`
	PermissionID int        `gorm:"not null" json:"-"`
	Permission   Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE;not null" json:"permission"`
}
