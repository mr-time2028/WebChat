package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID pgtype.UUID `gorm:"not null" json:"-"`
	Sender   User        `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;not null" json:"sender"`
	RoomID   pgtype.UUID `gorm:"not null" json:"-"`
	Room     Room        `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE;not null" json:"room"`
	Content  string      `gorm:"not null" json:"content"`
}
