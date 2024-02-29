package handlers

import (
	"github.com/joho/godotenv"
	"github.com/mr-time2028/WebChat/internal/config"
	"github.com/mr-time2028/WebChat/internal/database"
	"github.com/mr-time2028/WebChat/internal/models"
	"log"
	"os"
	"testing"
)

var testApp config.App

func addDefaultData() error {
	var defaultUsers = &models.User{
		Username: "defaultUser",
		Password: "defaultPass",
	}

	_, err := testApp.Models.User.InsertOneUser(defaultUsers)
	if err != nil {
		return err
	}

	return nil
}

func setUpTest() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("setUpTest error while load env file: ", err.Error())
	}

	// TODO: you can create a NewTestJWTAuth function to initial auth for the tests
	auth := &models.Auth{
		Secret: "TestSecretKey",
	}
	testApp.Auth = auth

	testDB, err := database.ConnectTestSQL()
	if err != nil {
		log.Fatal("setUpTest error while connect to the database: ", err.Error())
	}
	testApp.DB = testDB

	models.NewModels(testDB)
	testApp.Models = models.NewModelManager()

	err = testApp.DB.DropAllTables()
	if err != nil {
		log.Fatal("setUpTest error while drop all tables: ", err.Error())
	}

	err = models.AutoMigration()
	if err != nil {
		log.Fatal("setUpTest error while auto migration: ", err.Error())
	}

	err = addDefaultData()
	if err != nil {
		log.Fatal("setUpTest error while adding default user(s) data to the database: ", err)
	}

	NewHandlers(&testApp)
}

func tearDownTest() {
	_ = testApp.DB.DropAllTables()
}

func TestMain(m *testing.M) {
	setUpTest()
	exitCode := m.Run()
	tearDownTest()
	os.Exit(exitCode)
}
