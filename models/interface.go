package models

type UserInterface interface {
}

type RoomInterface interface {
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

func NewModels() *ModelManager {
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
