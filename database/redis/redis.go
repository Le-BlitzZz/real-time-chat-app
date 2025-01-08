package redis

import (
	"context"

	"github.com/Le-BlitzZz/real-time-chat-app/model/redis"
	go_redis "github.com/redis/go-redis/v9"
)

var nLastMessages int64 = 10

type RedisDb struct {
	*go_redis.Client
}

func New(address string) (*RedisDb, error) {
	redisDb := &RedisDb{go_redis.NewClient(&go_redis.Options{
		Addr: address,
	})}

	if err := redisDb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return redisDb, nil
}

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

func (db *RedisDb) AddUserToChat(ctx context.Context, userID uint, chatID uint) error {
	userChatsKey := UserChatsKey(userID)
	return db.SAdd(ctx, userChatsKey, chatID).Err()
}

func (db *RedisDb) GetUserChats(ctx context.Context, userID uint) ([]string, error) {
	userChatsKey := UserChatsKey(userID)
	return db.SMembers(ctx, userChatsKey).Result()
}

func (db *RedisDb) RemoveChat(ctx context.Context, chatID uint) error {
	chatUsersKey := ChatUsersKey(chatID)
	users, err := db.SMembers(ctx, chatUsersKey).Result()
	if err != nil {
		return err
	}

	chatKey := ChatKey(chatID)
	if err := db.Del(ctx, chatKey).Err(); err != nil {
		return err
	}

	for _, userIDstr := range users {
		userID := ParseUint(userIDstr)

		userChatsKey := UserChatsKey(userID)
		if err := db.SRem(ctx, userChatsKey, chatID).Err(); err != nil {
			return err
		}
	}

	return db.Del(ctx, chatUsersKey).Err()
}
