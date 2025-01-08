package sql

type FriendRequest struct {
	ID         uint `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	SenderID   uint `gorm:"not null;index"`
	ReceiverID uint `gorm:"not null;index"`

	Sender   User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	Receiver User `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE"`
}

type FriendRequestReceiver struct {
	Name string `json:"receiver_name" binding:"required"`
}

type FriendRequestAcceptReject struct {
	RequestID uint `json:"request_id" binding:"required"`
}
