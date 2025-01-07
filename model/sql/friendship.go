package sql

type FriendRequest struct {
	ID         uint `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	SenderID   uint `gorm:"not null;index"`
	ReceiverID uint `gorm:"not null;index"`

	Sender   User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	Receiver User `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE"`
}

type FriendRequestReceiverId struct {
	ReceiverID uint `json:"receiver_id" binding:"required"`
}

type FriendRequestAcceptReject struct {
	RequestID uint `json:"receiver_id" binding:"required"`
}