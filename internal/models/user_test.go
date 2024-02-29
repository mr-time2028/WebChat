package models

import (
	"github.com/google/uuid"
	"testing"
)

func TestUser_InsertOneUser(t *testing.T) {
	var testCases = []struct {
		name        string // name of the test
		user        *User  // user we want to insert to the database
		expectedErr bool   // do we expect any error from this query to the database?
	}{
		{
			"insert one user",
			&User{Username: "MrTime", Password: "MrTime1234"},
			false,
		},
		{
			"insert one user (duplicate username)",
			&User{Username: "Dav59", Password: "ABCD1234"}, // already added
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// test function
			user := User{}
			userID, err := user.InsertOneUser(tc.user)

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %v, err is %s", tc.expectedErr, err.Error())
			} else if err == nil && userID == uuid.Nil {
				t.Errorf("expected a notnil uuid user id, but it is nil uuid")
			}

			_ = testModelRepo.db.DropAllTables()
			_ = AutoMigration()
			_ = addDefaultData()
		})
	}
}

func TestUser_GetOneUser(t *testing.T) {
	var testCases = []struct {
		name        string // name of the test
		username    string // specific user id that we want to get user with it from database
		expectedErr bool   // do we expect any error from this query to the database?
	}{
		{
			"get user by username",
			"Dav59",
			false,
		},
		{
			"no rows",
			"SomeUsername",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// test function
			user := User{}
			_, err := user.GetUserByUsername(tc.username)

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %t, err is %s", tc.expectedErr, err.Error())
			}
		})

		_ = testModelRepo.db.DropAllTables()
		_ = AutoMigration()
		_ = addDefaultData()
	}
}

func TestUser_CheckIfExistsUser(t *testing.T) {
	var testCases = []struct {
		name        string
		username    string
		expectedErr bool
		isExists    bool
	}{
		{
			"user exists",
			"Dav59",
			false,
			true,
		},
		{
			"user does not exists",
			"MrTime",
			false,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := User{}
			isExists, err := user.CheckIfExistsUser(tc.username)

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %v, err is %s", tc.expectedErr, err.Error())
			} else if isExists != tc.isExists {
				t.Errorf("expected user exists %v but got user exists %v", tc.isExists, isExists)
			} else if isExists && err != nil {
				t.Errorf("got unknown error: %s", err.Error())
			}
		})

		_ = testModelRepo.db.DropAllTables()
		_ = AutoMigration()
		_ = addDefaultData()
	}
}
