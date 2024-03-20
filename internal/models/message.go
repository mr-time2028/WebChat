package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID string `gorm:"not null" json:"-"`
	Sender   User   `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;not null" json:"sender"`
	RoomID   string `gorm:"not null" json:"-"`
	Room     Room   `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE;not null" json:"room"`
	Content  string `gorm:"size:5000;not null" json:"content"`
}
