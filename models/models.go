package models

import (
	"github.com/mr-time2028/WebChat/database"
)

var db *database.DB

func RegisterModelsConfig(d *database.DB) {
	db = d
}
