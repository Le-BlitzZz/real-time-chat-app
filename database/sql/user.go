package sql

import (
	"errors"

	"github.com/Le-BlitzZz/real-time-chat-app/model/sql"
	"gorm.io/gorm"
)

func (db *SqlDb) CreateUser(user *sql.User) error {
	return db.DB.Create(user).Error
}

func (db *SqlDb) GetUserByName(name string) (*sql.User, error) {
	var user sql.User
	if err := db.Where("name = ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (db *SqlDb) GetUserByEmail(email string) (*sql.User, error) {
	var user sql.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err // Other errors
	}
	return &user, nil
}
