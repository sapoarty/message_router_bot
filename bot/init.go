package bot

import (
	"os"
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"message_router_bot/utils"
	"strconv"
)

var BotAPI *tgbotapi.BotAPI

func Init() {
	var err error
	log.Println(os.Getenv("BotToken"))
	BotAPI, err = tgbotapi.NewBotAPI(os.Getenv("BotToken"))
	if err != nil {
		panic(err)
	}
}

func InitUser(userID int) {
	utils.InitUserData(userID)
	handlerMap = GetHandlerMapForUser(userID)
	log.Println("InitUserData for user " + strconv.Itoa(userID))
}
