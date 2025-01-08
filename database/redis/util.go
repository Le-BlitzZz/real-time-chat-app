package redis

import (
	"fmt"
	"strconv"
)

func ChatKey(chatID uint) string {
	return fmt.Sprintf("chat:%d", chatID)
}

func UserChatsKey(userID uint) string {
	return fmt.Sprintf("user:%d:chats", userID)
}

func ChatUsersKey(chatID uint) string {
	return fmt.Sprintf("chat:%d:users", chatID)
}

func ParseUint(s string) uint {
	n64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(n64)
}
