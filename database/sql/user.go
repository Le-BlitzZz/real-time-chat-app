package sql

import (
	"errors"

	"github.com/Le-BlitzZz/real-time-chat-app/model/sql"
	"gorm.io/gorm"
)

func (db *SqlDb) CreateUser(user *sql.User) error {
	return db.Create(user).Error
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

func (db *SqlDb) GetFriends(userID uint) ([]sql.User, error) {
    var user sql.User
    if err := db.Preload("Friends").First(&user, userID).Error; err != nil {
        return nil, err
    }
	
	friends := make([]sql.User, len(user.Friends))
    for i, friend := range user.Friends {
        friends[i] = *friend
    }

    return friends, nil
}
