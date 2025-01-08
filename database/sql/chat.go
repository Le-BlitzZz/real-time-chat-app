package sql

import (
	"github.com/Le-BlitzZz/real-time-chat-app/model/sql"
)

func (db *SqlDb) CreateChat(chat *sql.Chat) error {
	return db.Create(chat).Error
}
