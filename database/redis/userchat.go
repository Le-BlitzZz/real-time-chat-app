package redis

import (
	"context"

	go_redis "github.com/redis/go-redis/v9"
)

func (db *RedisDb) CreateChat(ctx context.Context, chatID uint, meta map[string]interface{}) error {
	key := ChatMetaKey(chatID)
	return db.HSet(ctx, key, meta).Err()
}

func (db *RedisDb) AddUserToChat(ctx context.Context, userID uint, chatID uint) error {
	userChatsKey := UserChatsKey(userID)
	chatUsersKey := ChatUsersKey(chatID)

	if err := db.SAdd(ctx, userChatsKey, chatID).Err(); err != nil {
		return err
	}
	return db.SAdd(ctx, chatUsersKey, userID).Err()
}

func (db *RedisDb) RemoveUserFromChat(ctx context.Context, userID uint, chatID uint) error {
	userChatsKey := UserChatsKey(userID)
	chatUsersKey := ChatUsersKey(chatID)

	if err := db.SRem(ctx, userChatsKey, chatID).Err(); err != nil {
		return err
	}
	return db.SRem(ctx, chatUsersKey, userID).Err()
}

func (db *RedisDb) ChatExists(ctx context.Context, chatID uint) (bool, error) {
	chatMetaKey := ChatMetaKey(chatID)
	exists, err := db.Exists(ctx, chatMetaKey).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (db *RedisDb) GetUserChats(ctx context.Context, userID uint) ([]uint, error) {
	userChatsKey := UserChatsKey(userID)
	chatIDsStr, err := db.SMembers(ctx, userChatsKey).Result()
	if err != nil {
		return nil, err
	}

	var chatIDs []uint
	for _, idStr := range chatIDsStr {
		chatIDs = append(chatIDs, ParseUint(idStr))
	}
	return chatIDs, nil
}

func (db *RedisDb) RemoveChat(ctx context.Context, chatID uint) error {
	chatMetaKey := ChatMetaKey(chatID)
	chatUsersKey := ChatUsersKey(chatID)
	chatMessagesKey := ChatMessagesKey(chatID)

	users, err := db.SMembers(ctx, chatUsersKey).Result()
	if err != nil {
		return err
	}

	for _, userIDStr := range users {
		userID := ParseUint(userIDStr)
		userChatsKey := UserChatsKey(userID)
		if err := db.SRem(ctx, userChatsKey, chatID).Err(); err != nil {
			return err
		}
	}

	return db.Del(ctx, chatMetaKey, chatUsersKey, chatMessagesKey).Err()
}

func (db *RedisDb) GetOrCreatePrivateChatID(ctx context.Context, userID1, userID2 uint) (uint, error) {
	privateChatKey := PrivateChatKey(userID1, userID2)

	chatIDStr, err := db.Get(ctx, privateChatKey).Result()
	if err == go_redis.Nil {
		newChatID, err := GenerateChatID(db, ctx)
		if err != nil {
			return 0, err
		}

		if err := db.Set(ctx, privateChatKey, newChatID, 0).Err(); err != nil {
			return 0, err
		}

		return newChatID, nil
	} else if err != nil {
		return 0, err
	}

	return ParseUint(chatIDStr), nil
}

func (db *RedisDb) GetChatUsers(ctx context.Context, chatID uint) ([]uint, error) {
	chatUsersKey := ChatUsersKey(chatID)
	participantIDsStr, err := db.SMembers(ctx, chatUsersKey).Result()
	if err != nil {
		return nil, err
	}

	var participantIDs []uint
	for _, idStr := range participantIDsStr {
		participantIDs = append(participantIDs, ParseUint(idStr))
	}

	return participantIDs, nil
}
