package bot

import (
	"os"
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var BotAPI *tgbotapi.BotAPI

func InitBot() {
	var err error
	log.Println(os.Getenv("BotToken"))
	BotAPI, err = tgbotapi.NewBotAPI(os.Getenv("BotToken"))
	if err != nil {
		panic(err)
	}
}
