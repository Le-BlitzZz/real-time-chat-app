package sql

import "github.com/Le-BlitzZz/real-time-chat-app/model/sql"

func (db *SqlDb) CreateFriendRequest(senderID, receiverID uint) error {
	friendRequest := sql.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
	}
	return db.Create(&friendRequest).Error
}

func (db *SqlDb) GetFriendRequests(receiverID uint) ([]sql.FriendRequest, error) {
	var requests []sql.FriendRequest
	err := db.Preload("Sender").Preload("Receiver").
		Where("receiver_id = ?", receiverID).
		Find(&requests).Error
	return requests, err
}

func (db *SqlDb) AcceptFriendRequest(requestID uint) error {
	var request sql.FriendRequest
	if err := db.First(&request, requestID).Error; err != nil {
		return err
	}

	if err := db.Model(&sql.User{ID: request.SenderID}).Association("Friends").Append(&sql.User{ID: request.ReceiverID}); err != nil {
		return err
	}

	if err := db.Model(&sql.User{ID: request.ReceiverID}).Association("Friends").Append(&sql.User{ID: request.SenderID}); err != nil {
		return err
	}

	return db.Delete(&request).Error
}

func (db *SqlDb) RejectFriendRequest(requestID uint) error {
	return db.Delete(&sql.FriendRequest{}, requestID).Error
}
