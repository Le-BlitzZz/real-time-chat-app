package redis

import (
	"encoding/json"
	"time"
)

type Message struct {
	Username  string `json:"username"`
	Message   string `json:"message"`
	ChatID    uint   `json:"chat_id"`
	Timestamp int64  `json:"timestamp"` // Unix timestamp
}

type MessageIncoming struct {
	Message string `json:"message"`
}

func NewMessage(username, message string, chatID uint) Message {
	return Message{
		Username:  username,
		Message:   message,
		ChatID:    chatID,
		Timestamp: time.Now().Unix(),
	}
}

func (message *Message) Serialize() (string, error) {
	jsonMessage, err := json.Marshal(message)
	return string(jsonMessage), err
}

func DeserializeMessageIncoming(rawMessage string) (MessageIncoming, error) {
	var message MessageIncoming
	err := json.Unmarshal([]byte(rawMessage), &message)
	return message, err
}

func DeserializeMessage(rawMessage string) (Message, error) {
	var message Message
	err := json.Unmarshal([]byte(rawMessage), &message)
	return message, err
}

func DeserializeMessages(rawMessages []string) ([]Message, error) {
	var messages []Message
	for _, rawMessage := range rawMessages {
		message, err := DeserializeMessage(rawMessage)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
