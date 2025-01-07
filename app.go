package main

import (
	"fmt"
	"os"

	"github.com/Le-BlitzZz/real-time-chat-app/config"
	"github.com/Le-BlitzZz/real-time-chat-app/database"
	"github.com/Le-BlitzZz/real-time-chat-app/mode"
	"github.com/Le-BlitzZz/real-time-chat-app/router"
	"github.com/Le-BlitzZz/real-time-chat-app/runner"
)

var Mode = mode.LocalDev

func main() {
	mode.Set(Mode)

	fmt.Println("Starting real-time-chat-app")
	conf := config.Get()

	db, err := database.New(
		conf.Database.SQL.Connection,
		conf.Database.Redis.Addr,
		conf.DefaultUser.Name,
		conf.DefaultUser.Email,
		conf.DefaultUser.Password,
	)
	if err != nil {
		panic(err)
	}
	defer db.Cleanup()

	engine := router.Create(db, conf)

	if err := runner.Run(engine, conf); err != nil {
		fmt.Println("Server error: ", err)
		os.Exit(1)
	}
}
