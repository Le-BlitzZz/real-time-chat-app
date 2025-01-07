package redis

import "fmt"

func ChatKey(chatID string) string {
	return fmt.Sprintf("chat:%s", chatID)
}
