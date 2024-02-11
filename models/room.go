package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	ID       pgtype.UUID `json:"id"`
	Name     string      `json:"name"`
	Password string      `json:"password"`
}
