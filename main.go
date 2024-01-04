package main

import (
	"message_router_bot/database"
	"message_router_bot/config"
	"message_router_bot/bot"
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	log.Println("Бот запущен. Инициализация ...")
	config.Init()
	database.Init()
	bot.Init()
}

func main() {
	defer database.DB.Close()
	log.Println("Бот готов к работе, ждет команд.")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.BotAPI.GetUpdatesChan(u)
	bot.HandleUpdates(updates)
}