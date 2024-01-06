package bot

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"    
    "message_router_bot/constants"
    "message_router_bot/messages"
    "message_router_bot/database"
    "message_router_bot/structures"
    "message_router_bot/config"
    "message_router_bot/utils"
    "fmt"
    "log"
    "time"
    "strings"
    "os"
)

func HandleUpdates(updates tgbotapi.UpdatesChannel) {
    for update := range updates {
        // Обработка CallbackQuery от inline кнопок
        if update.CallbackQuery != nil {
            log.Printf("CallbackQuery received: %s from user id %d", update.CallbackQuery.Data, update.CallbackQuery.From.ID)
            // Получаем данные от CallbackQuery
            cq := update.CallbackQuery
            callBackData := cq.Data
            // callbackQueryID := cq.ID
            userID := int(cq.From.ID)
            chatID := cq.Message.Chat.ID
            // Если нет зарегистрированного состояния для этого пользователя, создаем новое
            if _, ok := config.UserStates[userID]; !ok {
                utils.InitUserData(userID)
            }

            // Определяем, на какую inline кнопку было совершено нажатие используя callBackData
            handleCategorySelection(callBackData, chatID, userID)
            // Создаем ответ на CallbackQuery
            callbackConfig := tgbotapi.NewCallback(cq.ID, "")
            if _, err := BotAPI.Request(callbackConfig); err != nil {
                log.Printf("Error answering callback query: %v", err)
                continue
            }
            // Проверяем, есть ли сообщение в обновлении
        } else if update.Message != nil {
            log.Printf(
                "Message [%s] from user id %d with id = %d (media_group = %s) (Caption = %s)", 
                update.Message.Text, 
                update.Message.From.ID, 
                update.Message.MessageID, 
                update.Message.MediaGroupID,
                update.Message.Caption,
            )
            message := update.Message
            chatID := message.Chat.ID
            userID := int(message.From.ID)
            text := strings.TrimSpace(message.Text)
            
            if _, ok := config.UserStates[userID]; !ok {
                utils.InitUserData(userID)
                handlerMap = GetHandlerMapForUser(userID)
            }

            if message.Photo != nil && message.MediaGroupID != "" {
                log.Printf("Msg with photo group %s", message.MediaGroupID)
                groupID := message.MediaGroupID
                if _, ok := structures.MediaGroups[userID]; !ok {
                    structures.InitMediaGroupsForUser(userID)
                }
                if _, ok := structures.MediaGroups[userID].Groups[groupID]; !ok {
                    structures.InitMediaGroupsForGroupID(userID, groupID)
                    log.Println("Устанавливаем таймер, так как создана новая группа")
                    structures.MediaGroups[userID].Groups[groupID].Timer = time.AfterFunc(constants.MediaGroupWaitTime, func() {
                        log.Println("Таймер сработал для медиа-группы")
                        ForwardMessage(structures.MediaGroups[userID].Groups[groupID].Messages)
                        structures.ClearMediaGroupMessages(userID, groupID)
                    })
                }
                structures.AddMediaGroupsMessages(userID, groupID, message)
                log.Printf("photo in group = %d", len(structures.MediaGroups[userID].Groups[groupID].Messages))
                // Не отправляем сразу - сначала накапливаем всю группу
                if IsMediaGroupComplete(userID, groupID) {
                    log.Println("IsMediaGroupComplete = true")
                    ForwardMessage(structures.MediaGroups[userID].Groups[groupID].Messages)
                    structures.ClearMediaGroupMessages(userID, groupID)
                }
                continue
            }

            chatType := message.Chat.Type
            // Если нет зарегистрированного состояния для этого пользователя, создаем новое
            botCommand, commandExists := handlerMap[text] 

            switch {
            case chatType == "private":
                // Если сообщение отправлено боту
                handleBotMessage(message)

            case utils.IsGroupDefault(chatID, userID):
                // Обработка сообщений, если группа - по-умолчанию
                handleDefaultGroupMessage(message)

            case commandExists && !config.UserStates[userID].ExpectingInput:
                // Если команда существует и пользователь не ожидает ввода
                if botCommand.NeedReply {
                    // Если команда требует ввода, обновляем состояние пользователя и скрываем меню
                    askToReply(message)
                } else {
                    // Иначе вызываем обработчик команды
                    botCommand.Handler(text, chatID, userID)
                }

            case config.UserStates[userID].ExpectingInput:
                // Если пользователь ожидает ввода, перенаправляем его функции handleUserInput
                handleUserInput(text, chatID, userID, config.UserStates[userID].UserCommand)

            case message.LeftChatMember != nil:
                leftMember := update.Message.LeftChatMember
                // Если ID ушедшего пользователя совпадает с ID бота, то бот был удален из группы
                if leftMember.ID == BotAPI.Self.ID {
                    log.Printf("Bot was removed from group: %d", chatID)
                    database.DeleteChatData(chatID, userID)
                }
            case message.NewChatMembers != nil: // Бот добавлен в группу
                // Перебираем всех новых участников группы
                for _, newUser := range update.Message.NewChatMembers {
                    if newUser.ID == BotAPI.Self.ID {
                        // Бот был добавлен в группу
                        handleStart("", chatID, userID)
                    }
                }
            }
        }
    }
}

func handleUserInput(input string, chatID int64, userID int, userCommand string){
    log.Println("userCommand" + userCommand + "Lang " + config.UserStates[userID].Lang)
    handlerMap[userCommand].Handler(input, chatID, userID)
    // В данном подходе мы сбрасываем состояние пользователя после обработки ввода
    expectingInput := false
    utils.SetUserData(userID, &expectingInput, nil, nil)
    msg := fmt.Sprintf(
        messages.AwaitingCommand[config.UserStates[userID].Lang], 
        os.Getenv("BotLink"),
    )
    sendMessage(msg, chatID, GetMenu(userID), "Markdown")
}


func handleCategorySelection(category string, chatID int64, userID int) {
    // Получаем карту категорий.
    categoriesMap, err := GetDefaultCategories(userID)
    if err != nil {
        log.Panic(err)
        return
    }

    // Проверяем, есть ли выбранная категория.
    words, exists := categoriesMap[category]
    if !exists {
        log.Printf("Категория %s не найдена", category)
        return
    }

    categoryWordsString := strings.Join(words, ",")
    handleAddKeyWordsForGroup(categoryWordsString, chatID, userID)
}

func handleDefaultGroupMessage(message *tgbotapi.Message) {
    text := strings.TrimSpace(message.Text)
    chatID := message.Chat.ID
    userID := int(message.From.ID)
    userLang := config.UserStates[userID].Lang

    if (text == messages.CommandStart[userLang] || text == messages.CommandHelp[userLang]) {
        // Если приветственное сообщение или нужна помощь, выводим текст и короткое меню
        msg := fmt.Sprintf(
            messages.GreetingsMessage[userLang] +
            messages.GroupIsAlreadyDefaultGroup[userLang] + "\n" +
            messages.AwaitingCommand[userLang], 
            os.Getenv("BotLink"),
        )
        sendMessage(msg, chatID, GetShortMenu(userID), "Markdown")
    } else if (text == messages.CommandPrintAllKeywords[userLang]) {
        handlePrintAllKeywords(text, chatID, userID)
    } else if (text == messages.CommandChangeLang[userLang]) {
        // Смена языка
        var userLang string
        if (config.UserStates[userID].Lang == "ru") {
            userLang = "en"
        } else {
            userLang = "ru"
        }
        utils.SetUserData(userID, nil, nil, &userLang)
        handlerMap = GetHandlerMapForUser(userID)

        sendMessage(messages.LangIsChanged[config.UserStates[userID].Lang], chatID, GetShortMenu(userID))
    }
}

func handleBotMessage(message *tgbotapi.Message) {
    text := strings.TrimSpace(message.Text)
    chatID := message.Chat.ID
    userID := int(message.From.ID)
    userLang := config.UserStates[userID].Lang

    if (text == messages.CommandStart[userLang] || text == messages.CommandHelp[userLang]) {
        // Если приветственное сообщение или нужна помощь, выводим текст и короткое меню
        msg := fmt.Sprintf(
            messages.GreetingsMessage[config.UserStates[userID].Lang] + messages.AwaitingCommand[config.UserStates[userID].Lang], 
            os.Getenv("BotLink"),
        )
        sendMessage(msg, chatID, GetShortMenu(userID), "Markdown")
    } else if (text == messages.CommandChangeLang[userLang]) {
        // Смена языка
        var userLang string
        if (config.UserStates[userID].Lang == "ru") {
            userLang = "en"
        } else {
            userLang = "ru"
        }
        utils.SetUserData(userID, nil, nil, &userLang)
        handlerMap = GetHandlerMapForUser(userID)

        sendMessage(messages.LangIsChanged[config.UserStates[userID].Lang], chatID, GetShortMenu(userID))
    } else if (text == messages.CommandPrintAllKeywords[userLang]) {
        handlePrintAllKeywords(text, chatID, userID)
    } else {
        // Иначе просто пересылаем сообщение, куда нужно и удаляем исходное
        log.Println("ForwardMessage")
        ForwardMessage([]*tgbotapi.Message{message})
        log.Println("Delete message ...")
        if message != nil {
            log.Println(" message != nil")
            msg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
            if _, err := BotAPI.Request(msg); err != nil {
                log.Println("Error deleting message:", err)
            }
            log.Println(" message deleted")
        }    
    }
}

// Здесь проверяем, завершена ли группа медиафайлов
func IsMediaGroupComplete(userID int, groupID string) bool {
    userGroups, userExists := structures.MediaGroups[userID]
    if (!userExists) {
        return false
    }
    mediaGroup, groupExists := userGroups.Groups[groupID]
    if (!groupExists) {
        return false
    }
    if len(mediaGroup.Messages) >= constants.MaxMediaGroupMessages {
        if mediaGroup.Timer != nil {
            mediaGroup.Timer.Stop()
            mediaGroup.Timer = nil
        }
        return true
    }

    return false
}

