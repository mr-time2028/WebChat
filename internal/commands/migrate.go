package commands

import (
	"fmt"
	"github.com/mr-time2028/WebChat/internal/database"
	"github.com/mr-time2028/WebChat/internal/models"
	"log"
	"reflect"
)

// migrate is a simple command to create a superuser in database
func migrate() {
	db, err := database.ConnectSQL()
	if err != nil {
		log.Fatal("cannot connect to the database: ", err)
	}

	models := models.NewModelManager()
	modelsValue := reflect.ValueOf(*models)
	for i := 0; i < modelsValue.NumField(); i++ {
		field := modelsValue.Field(i)
		if field.Kind() == reflect.Interface {
			model := field.Interface()
			if err = db.GormDB.AutoMigrate(model); err != nil {
				log.Fatal("error while migration: ", err)
			}
		}
	}

	fmt.Println("migration was successful!")
}
