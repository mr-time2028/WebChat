package models

import (
	"log"
	"reflect"
)

type UserInterface interface {
	InsertOneUser(user *User) (string, error)
	CheckIfExistsUser(username string) (bool, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id string) (*User, error)
}

type RoomInterface interface {
	GetAllRooms() ([]*Room, error)
	InsertOneRoom(room *Room) (string, error)
}

type MessageInterface interface {
}

type PermissionInterface interface {
}

type UserRoomInterface interface {
}

type MessageReceiverInterface interface {
}

type UserRoomPermissionInterface interface {
}

type ModelManager struct {
	User               UserInterface
	Room               RoomInterface
	Message            MessageInterface
	Permission         PermissionInterface
	UserRoom           UserRoomInterface
	MessageReceiver    MessageReceiverInterface
	UserRoomPermission UserRoomPermissionInterface
}

func NewModelManager() *ModelManager {
	return &ModelManager{
		User:               &User{},
		Room:               &Room{},
		Message:            &Message{},
		Permission:         &Permission{},
		UserRoom:           &UserRoom{},
		MessageReceiver:    &MessageReceiver{},
		UserRoomPermission: &UserRoomPermission{},
	}
}

func AutoMigration() error {
	models := NewModelManager()
	modelsValue := reflect.ValueOf(*models)
	for i := 0; i < modelsValue.NumField(); i++ {
		field := modelsValue.Field(i)
		if field.Kind() == reflect.Interface {
			model := field.Interface()
			if err := ModelRepo.db.GormDB.AutoMigrate(model); err != nil {
				log.Fatal("error while migration: ", err)
			}
		}
	}

	return nil
}
