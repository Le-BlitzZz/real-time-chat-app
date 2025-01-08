package sql

import (
	"github.com/Le-BlitzZz/real-time-chat-app/model/sql"
)

func (db *SqlDb) CreateChat(chat *sql.Chat) error {
	return db.Create(chat).Error
}

func (db *SqlDb) GetChatByID(chatID uint) (*sql.Chat, error) {
	var chat sql.Chat
	if err := db.DB.First(&chat, chatID).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}
