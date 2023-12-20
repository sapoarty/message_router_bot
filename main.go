	package main

	import (
		"message_router_bot/database"
		"message_router_bot/bot"
		"strings"
		"log"
		"message_router_bot/config"
		"github.com/go-telegram-bot-api/telegram-bot-api"
		// _ "github.com/mattn/go-sqlite3"
	)

	type BotMessagesHandler func(text string, chatID int64, userID int)

	var handlerMap = map[string]BotMessagesHandler{
		"/add_keywords" : handleAddKeyWordsForChannel,
		"/delete_keywords" : handleDeleteKeyWordsForChannel,
	}

	func main() {
		log.Println("Бот запущен. Инициализация ...")
		database.InitDb()
		defer database.DB.Close()
		config.Init()
		bot.InitBot()
		log.Println("Бот готов к работе, ждет команд.")

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, _ := bot.BotAPI.GetUpdatesChan(u)
		handleUpdates(updates)
	}

	func handleUpdates(updates tgbotapi.UpdatesChannel) {
		for update := range updates {
			if update.Message != nil {
				log.Printf("Message from user id %d", update.Message.From.ID)
				message := update.Message
				chatID := message.Chat.ID
				userID := message.From.ID
				splitMessage := strings.Fields(message.Text)
				command := splitMessage[0]

				if handler, exists := handlerMap[command]; exists {
					text := ""
					if len(splitMessage) > 1 {
						text = strings.Join(splitMessage[1:], " ")
					}
					handler(text, chatID, userID)
				} else {
					bot.ForwardMessage(message)
				}
			} else {
				continue
			}
		}
	}

	func handleAddKeyWordsForChannel(text string, chatID int64, userID int) {
		keywordsList := strings.Split(text, ",")

		if len(keywordsList) == 0 {
			bot.BotAPI.Send(tgbotapi.NewMessage(chatID, "Keywords list is empty, please add using /set_kwlist command and ',' separator (without spaces!)"))
			return
		}

		database.AddKeywords(keywordsList, chatID, userID)
		log.Println("handleAddKeyWordsForChannel.chatID %d userID %d", chatID, userID)
		printKeywordsToChannel(chatID, userID)
	}

	func handleDeleteKeyWordsForChannel(text string, chatID int64, userID int) {
		keywordsList := strings.Split(text, ",")

		if len(keywordsList) == 0 {
			bot.BotAPI.Send(tgbotapi.NewMessage(chatID, "Keywords list is empty, please add using /delete_keywords command and ',' separator (without spaces!)"))
			return
		}

		database.DeleteKeywords(keywordsList, chatID, userID)
		log.Println("handleDeleteKeyWordsForChannel.chatID %d", chatID)
		printKeywordsToChannel(chatID, userID)
	}

	func printKeywordsToChannel(chatID int64, userID int) {
		err, channelName := bot.GetChannelNameByChatId(chatID)
		if err != nil {
			log.Printf("Err GetChannelNameByChatId: %s\n", err)
			return
		}

		log.Printf("printKeywordsToChannel %s", channelName)
		keywords := database.GetKeywordsForUserChatID(chatID, userID)
		if err != nil {
			log.Printf("Err GetKeywordsByChannelName: %s\n", err)
			return
		}
		if len(keywords) == 0 {
			bot.BotAPI.Send(tgbotapi.NewMessage(chatID, "Keywords list is empty, please add using /set_kwlist command and ',' separator (without spaces!)"))
			return
		}
		keywordsText := strings.Join(keywords, ", ")
		message := tgbotapi.NewMessage(
			chatID, 
			"Keyword list for channel '" + channelName + "': " + keywordsText,
		)

		_, err = bot.BotAPI.Send(message)
		if err != nil {
			log.Printf("%q: %s\n", err)
			return
		}
	}
