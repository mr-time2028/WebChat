package models

import (
	"github.com/mr-time2028/WebChat/database"
)

var ModelRepo *ModelRepository

type ModelRepository struct {
	db *database.DB
}

func NewModelsRepository(d *database.DB) *ModelRepository {
	return &ModelRepository{
		db: d,
	}
}

func NewModels(r *ModelRepository) {
	ModelRepo = r
}
