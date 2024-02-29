package models

import (
	"github.com/joho/godotenv"
	"github.com/mr-time2028/WebChat/internal/database"
	"log"
	"os"
	"testing"
)

var testModelRepo ModelRepository

func addDefaultData() error {
	var defaultUsers = []*User{
		{Username: "Dav59", Password: "DavidPass"},
	}

	if err := testModelRepo.db.GormDB.CreateInBatches(defaultUsers, len(defaultUsers)).Error; err != nil {
		return err
	}

	return nil
}

func setUpTest() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("setUpTest error while load env file: ", err.Error())
	}

	testDB, err := database.ConnectTestSQL()
	if err != nil {
		log.Fatal("setUpTest error while connect to the database: ", err.Error())
	}
	testModelRepo.db = testDB

	NewModels(testDB)

	err = AutoMigration()
	if err != nil {
		log.Fatal("setUpTest error while auto migration: ", err.Error())
	}

	err = addDefaultData()
	if err != nil {
		log.Fatal("setUpTest error while adding default user(s) data to the database: ", err.Error())
	}
}

func tearDownTest() {
	_ = testModelRepo.db.DropAllTables()
}

func TestMain(m *testing.M) {
	setUpTest()
	exitCode := m.Run()
	tearDownTest()
	os.Exit(exitCode)
}
