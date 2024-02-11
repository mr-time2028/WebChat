package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       pgtype.UUID `json:"id"`
	Username string      `json:"username"`
	Password string      `json:"password"`
}
