package redis

import (
	"context"

	"github.com/Le-BlitzZz/real-time-chat-app/model/redis"
)

func (db *RedisDb) PublishMessage(ctx context.Context, message redis.Message) error {
	jsonMessage, err := message.Serialize()
	if err != nil {
		return err
	}

	return db.Publish(ctx, ChatKey(message.ChatID), jsonMessage).Err()
}

func (db *RedisDb) SaveMessage(ctx context.Context, message redis.Message) error {
	jsonMessage, err := message.Serialize()
	if err != nil {
		return err
	}

	chat := ChatKey(message.ChatID)
	if err := db.RPush(ctx, chat, jsonMessage).Err(); err != nil {
		return err
	}

	return db.LTrim(ctx, chat, -nLastMessages, -1).Err()
}

func (db *RedisDb) GetRecentMessages(ctx context.Context, chatID uint) ([]redis.Message, error) {
	chat := ChatKey(chatID)
	rawMessages, err := db.LRange(ctx, chat, -nLastMessages, -1).Result()
	if err != nil {
		return nil, err
	}

	return redis.DeserializeMessages(rawMessages)
}
