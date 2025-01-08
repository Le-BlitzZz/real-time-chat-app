package redis

import "fmt"

func ChatKey(chatID uint) string {
	return fmt.Sprintf("chat:%s", chatID)
}

func UserChatsKey(userID uint) string {
	return fmt.Sprintf("user:%s:chats", userID)
}

func ChatUsersKey(chatID uint) string {
	return fmt.Sprintf("chat:%s:users", chatID)
}
