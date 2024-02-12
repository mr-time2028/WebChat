package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	Identifier string    `gorm:"not null" json:"identifier"`
	Password   string    `gorm:"not null" json:"password"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) error {
	// set unique identifier
	uuidVal, _ := uuid.NewRandom()
	r.Identifier = uuidVal.String()

	// set password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	r.Password = string(hashedPassword)

	return nil
}
