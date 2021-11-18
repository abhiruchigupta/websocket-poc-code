package ws

import "time"

type WsMessage interface {
	GetUserID() string
	GetStoreID() int64
	GetMessage() string
	GetSenderID() string
	GetMessageType() string
}

type DbMessage struct {
	ID         *int64
	UserID     *string
	SenderID   *string
	Message    *string
	ReceivedAt *time.Time
	Sent       *bool
}

type MessageStore interface {
	Initialize() error
	StoreMessage(userID string, message string, receivedAt time.Time) (int64, error)
	RetrieveUnsentMessages(userID string) []DbMessage
	ConfirmSentMessage(id int64) error
	Close()
}
