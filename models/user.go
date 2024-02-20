package models

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username string    `gorm:"size:30;uniqueIndex;not null" json:"username"`
	IsActive bool      `gorm:"not null" json:"is_active"`
	Password string    `gorm:"size:60;not null" json:"password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}

// InsertOneUser insert one user to the database
func (u *User) InsertOneUser(user *User) (uuid.UUID, error) {
	result := ModelRepo.db.GormDB.Create(user)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return user.ID, nil
}

// CheckIfExistsUser check if user already exists in the database
func (u *User) CheckIfExistsUser(username string) (bool, error) {
	var user *User
	condition := User{Username: username}
	result := ModelRepo.db.GormDB.Where(condition).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (u *User) GetUserByUsername(username string) (*User, error) {
	var user *User
	condition := User{Username: username}
	result := ModelRepo.db.GormDB.Where(condition).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
