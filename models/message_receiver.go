package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageReceiver struct {
	gorm.Model
	ReceiverID uuid.UUID `gorm:"not null" json:"-"`
	Receiver   User      `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE;not null" json:"user"`
	MessageID  int       `gorm:"not null" json:"-"`
	Message    Message   `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE;not null" json:"message"`
}
