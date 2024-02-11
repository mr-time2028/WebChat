package models

type UserInterface interface {
}

type RoomInterface interface {
}

type PermissionInterface interface {
}

type UserRoomInterface interface {
}

type UserRoomPermissionInterface interface {
}

type UserRoomRoleInterface interface {
}

type ModelManager struct {
	User               UserInterface
	Room               RoomInterface
	Permission         PermissionInterface
	UserRoom           UserRoomInterface
	UserRoomPermission UserRoomPermissionInterface
	UserRoomRole       UserRoomRoleInterface
}

func NewModels() *ModelManager {
	return &ModelManager{
		User:               &User{},
		Room:               &Room{},
		Permission:         &Permission{},
		UserRoom:           &UserRoom{},
		UserRoomPermission: &UserRoomPermission{},
		UserRoomRole:       &UserRoomRole{},
	}
}
