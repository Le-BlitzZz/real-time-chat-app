package sql

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/Le-BlitzZz/real-time-chat-app/auth/password"
	"github.com/Le-BlitzZz/real-time-chat-app/model/sql"
)

const createDefaultUserIfNotExist = true

type SqlDb struct {
	*gorm.DB
}

func (db *SqlDb) Close() {
	sqlDb, _ := db.DB.DB()
	sqlDb.Close()
}

func New(sqlConnection, defaultUserName, defaultUserEmail, defaultPassword string) (*SqlDb, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: sqlConnection,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// We normally don't need that much connections, so we limit them. F.ex. mysql complains about
	// "too many connections", while load testing Gotify.
	sqlDb, _ := db.DB()
	sqlDb.SetMaxOpenConns(10)

	// Mysql has a setting called wait_timeout, which defines the duration
	// after which a connection may not be used anymore.
	// The default for this setting on mariadb is 10 minutes.
	// See https://github.com/docker-library/mariadb/issues/113
	sqlDb.SetConnMaxLifetime(9 * time.Minute)

	if err := db.AutoMigrate(new(sql.User)); err != nil {
		return nil, err
	}


	var userCount int64 = 0
	db.Find(new(sql.User)).Count(&userCount)
	if createDefaultUserIfNotExist && userCount == 0 {
		db.Create(&sql.User{
			Name: defaultUserName,
			Email: defaultUserEmail,
			Password: password.CreatePassword(defaultPassword),
			Admin: true,
		})
	}

	return &SqlDb{db}, nil
}
