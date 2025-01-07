package database

import (
	"github.com/Le-BlitzZz/real-time-chat-app/database/redis"
	"github.com/Le-BlitzZz/real-time-chat-app/database/sql"
	_ "github.com/redis/go-redis/v9"
)

type Database struct {
	SQL   *sql.SqlDb
	Redis *redis.RedisDb
}

func New(sqlConnection, redisAddr, defaultUserName, defaultUserEmail, defaultPassword string) (*Database, error) {
	sqlDb, err := sql.New(sqlConnection, defaultUserName, defaultUserEmail, defaultPassword)
	if err != nil {
		return nil, err
	}

	redisDb, err := redis.New(redisAddr)
	if err != nil {
		return nil, err
	}

	return &Database{
		SQL:   sqlDb,
		Redis: redisDb,
	}, nil
}

func (db *Database) Cleanup() {
	db.SQL.Close()
	db.Redis.Close()
}
