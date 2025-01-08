package redis

import (
	"context"
	"fmt"
	"strconv"
)

var chatIDCounterKey = "chat:counter" // Redis key for the global chat counter.
const GlobalChatID uint = 0

func GenerateChatID(db *RedisDb, ctx context.Context) (uint, error) {
	chatID, err := db.Incr(ctx, chatIDCounterKey).Result()
	if err != nil {
		return 0, err
	}
	return uint(chatID), nil
}

func PrivateChatKey(userID1, userID2 uint) string {
	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}
	return fmt.Sprintf("private:%d:%d", userID1, userID2)
}

func ChatMetaKey(chatID uint) string {
	return fmt.Sprintf("chat:%d:meta", chatID)
}

func ChatUsersKey(chatID uint) string {
	return fmt.Sprintf("chat:%d:users", chatID)
}

func ChatMessagesKey(chatID uint) string {
	return fmt.Sprintf("chat:%d:messages", chatID)
}

func UserChatsKey(userID uint) string {
	return fmt.Sprintf("user:%d:chats", userID)
}

func ParseUint(s string) uint {
	n64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(n64)
}
