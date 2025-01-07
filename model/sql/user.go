package sql

// User represents the user entity stored in the database.
// It is defined by either Name or Email, meaning users can log in and be identified using either.
type User struct {
	ID       uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Name     string `gorm:"type:varchar(180);unique_index"`
	Email    string `gorm:"type:varchar(180);unique_index"`
	Password []byte
	Admin    bool
	Friends  []*User `gorm:"many2many:user_friends"`
}

// Used for creating a new user.
type UserCreation struct {
	Name     string `binding:"required" json:"name" query:"name" form:"name"`
	Email    string `binding:"required" json:"email" query:"email" form:"email"`
	Password string `binding:"required" json:"password" query:"password" form:"password"`
	Admin    bool   `json:"admin" form:"admin" query:"admin"`
}
