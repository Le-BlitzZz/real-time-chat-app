package redis

import "context"

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

func (db *RedisDb) RemoveUserFromChat(ctx context.Context, userID uint, chatID uint) error {
	userChatsKey := UserChatsKey(userID)
	chatUsersKey := ChatUsersKey(chatID)

	if err := db.SRem(ctx, userChatsKey, chatID).Err(); err != nil {
		return err
	}

	return db.SRem(ctx, chatUsersKey, userID).Err()
}
