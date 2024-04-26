package models

import (
	"github.com/mr-time2028/WebChat/internal/database"
)

var ModelRepo *ModelRepository

type ModelRepository struct {
	db *database.DB
}

func NewModels(d *database.DB) {
	ModelRepo = &ModelRepository{
		db: d,
	}
}
