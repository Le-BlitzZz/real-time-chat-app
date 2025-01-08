package sql

const GlobalChatID = 1

const (
	TypeGlobal  = "global"
	TypePrivate = "private"
	TypeGroup   = "group"
)

type Chat struct {
	ID   uint    `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Name *string `gorm:"type:varchar(255);null"`
	Type string  `gorm:"type:enum('global','private','group');not null"`
}

var GlobalChat = &Chat{ID: 1, Type: TypeGlobal}
var PrivateChat = &Chat{Type: TypePrivate}

func NewGroupChat(name string) *Chat {
	return &Chat{
		Name: &name,
		Type: TypeGroup,
	}
}
