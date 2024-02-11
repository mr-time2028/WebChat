package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type UserRoom struct {
	gorm.Model
	UserID pgtype.UUID `json:"user_id"`
	RoomID pgtype.UUID `json:"room_id"`
}
