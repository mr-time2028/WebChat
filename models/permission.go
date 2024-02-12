package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name        string `gorm:"size:255;uniqueIndex;not null" json:"name"`
	Description string `json:"size:1024;description"`
}
