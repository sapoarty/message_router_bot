package bot

import (
    "message_router_bot/constants"
    "message_router_bot/messages"
    "message_router_bot/structures"
    "message_router_bot/config"
    "message_router_bot/utils"
    "message_router_bot/database"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "fmt"
    "strings"
    "time"
)

var handlerMap map[string]structures.BotMessagesHandler

func ForwardMessage(message *tgbotapi.Message) error {
    var matchFound = false
    var usedKeywordsForChatsList = make(map[int64]string)
    userID := message.From.ID
    botChatID := message.Chat.ID
    messageToDelete := tgbotapi.NewDeleteMessage(botChatID, message.MessageID)
    userLang := config.UserStates[userID].Lang
    text := ""
    if message.Text != "" { // Проверяем, что сообщение содержит текст
        text = message.Text
        log.Println(text)
    }

    // Итерируемся через наши ключевые слова и их идентификаторы чата
    for keyword, chatID := range structures.UsersKeywordsChatsMap[userID] {
        if _, ok := usedKeywordsForChatsList[chatID]; ok {
            continue
        }
        // Если сообщение содержит ключевое слово (без учета регистра)
        if text != "" && strings.Contains(strings.ToLower(text), strings.ToLower(keyword)) {
            err, groupName := GetGroupNameByChatId(chatID)
            if err != nil {
                log.Printf("Err GetGroupNameByChatId: %s\n", err)
                return err
            }
            log.Println("Contains keyword:", keyword, "Group Name:", groupName)

            // Создаем запрос на пересылку сообщения
            copyMsg := tgbotapi.NewForward(chatID, botChatID, message.MessageID)
            _, err = BotAPI.Send(copyMsg)
            if err != nil {
                if (strings.Contains(err.Error(), "Forbidden: the group chat was deleted")) {
                    _, groupName := GetGroupNameByChatId(chatID)
                    msg := fmt.Sprintf(messages.ChatDeleted[userLang], groupName, keyword)
                    sendMessage(msg, botChatID)
                    database.DeleteChatData(chatID, userID)
                } else {
                    log.Println("Error during forwarding message:", err)
                    return err
                }
            } else {
                matchFound = true
                usedKeywordsForChatsList[chatID] = keyword
            }
        }
    }

    if !matchFound {
        // Если в сообщении не нашлось ключевого слова, отправляем в группу по-умолчанию
        chatID := structures.GetUserChatIDForKeyword(constants.DefaultGroup, userID)
        if chatID == 0 {
            sendErrorMessage := tgbotapi.NewMessage(botChatID, messages.DefaulGroupIsNotSet[userLang])
            BotAPI.Send(sendErrorMessage)
            return nil
        }
        log.Println("Default Group")

        // Создаем запрос на пересылку сообщения в группу по-умолчанию
        copyMsg := tgbotapi.NewForward(chatID, botChatID, message.MessageID)
        _, err := BotAPI.Send(copyMsg)
        if err != nil {
            log.Println("Error during forwarding message to default group:", err)
            return err
        }
    }

    if (matchFound) {
        var chatNamesListStr string
        for chat, keyword := range usedKeywordsForChatsList {
            _, chatName := GetGroupNameByChatId(chat)
            chatNamesListStr += fmt.Sprintf("%s(%s), ", chatName,keyword)
        }
        chatNamesListStr = strings.TrimSuffix(chatNamesListStr, ", ")
        runes := []rune(text)
        if len(runes) > 100 {
            text = string(runes[:100]) + "..."
        }
        msg := fmt.Sprintf(messages.MessageHasBeenForwardedToChatsUsingKeywords[userLang], text, chatNamesListStr)
        message, err := sendMessage(msg, botChatID)
        if err != nil {
            return err
        }
        go countdownAndDeleteMessage(botChatID, message)
    }

    BotAPI.DeleteMessage(messageToDelete)
    return nil
}


func askToReply(message *tgbotapi.Message) {
    text := strings.TrimSpace(message.Text)
    chatID := message.Chat.ID
    userID := message.From.ID
    userLang := config.UserStates[userID].Lang

    expectingInput := true
    utils.SetUserData(userID, &expectingInput, &text, nil)

    sendMessage(
        messages.InputKeywordsListRequest[userLang], 
        chatID,
        tgbotapi.ReplyKeyboardRemove{
            RemoveKeyboard: true,
            Selective:      false,
        },
    )
}

func printKeywordsToGroup(chatID int64, userID int) {
    userLang := config.UserStates[userID].Lang
    err, groupName := GetGroupNameByChatId(chatID)
    if err != nil {
        log.Printf("Err GetGroupNameByChatId: %s\n", err)
        return
    }

    log.Printf("printKeywordsToGroup %s", groupName)
    keywords := structures.GetKeywordsForUserChatID(chatID, userID)
    if err != nil {
        log.Printf("Err GetKeywordsByGroupName: %s\n", err)
        return
    }

    var msgText string
    if len(keywords) == 0 {
        msgText = messages.KeywordsListEmpty[userLang]
    } else {
        msgText = fmt.Sprintf(messages.KeywordsListForGroup[userLang], groupName, strings.Join(keywords, ", "))
    }
    sendMessage(msgText, chatID)
}

func GetGroupNameByChatId(chatID int64) (error, string) {
    chat, err := BotAPI.GetChat(tgbotapi.ChatConfig{ChatID: chatID})
    if err != nil {
        log.Printf("GetGroupNameByChatId err %s", err)
        return err, ""
    }
    return nil, chat.Title
}

func sendMessage(msgText string, chatID int64, additionalArgs ...interface{}) (tgbotapi.Message, error) {
    message := tgbotapi.NewMessage(chatID, msgText)

    // Обрабатываем дополнительные аргументы, если они есть
    for _, arg := range additionalArgs {
        switch argTyped := arg.(type) {
        case tgbotapi.InlineKeyboardMarkup:
            message.ReplyMarkup = argTyped
        case tgbotapi.ReplyKeyboardMarkup:
            message.ReplyMarkup = argTyped
        case tgbotapi.ReplyKeyboardRemove:
            message.ReplyMarkup = argTyped
        case string:
            // Используем константы из пакета tgbotapi для ParseMode
            if argTyped == "Markdown" || argTyped == "MarkdownV2" {
                message.ParseMode = tgbotapi.ModeMarkdown
            } else if argTyped == "HTML" {
                message.ParseMode = tgbotapi.ModeHTML
            }
        default:
            // Обрабатываем другие случаи, например если требуется установить другие поля сообщения
            log.Printf("Неизвестный или неподдерживаемый тип для дополнительных аргументов: %T", argTyped)
        }
    }


    // Отправляем сообщение
    result, err := BotAPI.Send(message)
    if err != nil {
        log.Printf("Ошибка при отправке сообщения в чат ID %d: %v", chatID, err)
        return result, err
    }

    return result, nil
}

func setGroupDesc(chatID int64, userID int) {
    userLang := config.UserStates[userID].Lang
    log.Printf("setGroupDesc chatID %d, userID %d", chatID, userID)

    var descText string
    if (utils.IsGroupDefault(chatID, userID)) {
        descText = fmt.Sprintf(messages.DefaultGroupDesc[userLang])
    } else {
        keywords := structures.GetKeywordsForUserChatID(chatID, userID)
        if (len(keywords) > 0) {
            descText = fmt.Sprintf(
                messages.GroupDesc[userLang],
                strings.Join(keywords, ", "),
            )
        }
    }

    _, err := BotAPI.Send(tgbotapi.SetChatDescriptionConfig{
        ChatID:      chatID,
        Description: descText,
    })
    if err != nil {
        log.Printf("Error setting group description for chatID %d: %v", chatID, err)
    }
}

func countdownAndDeleteMessage(chatID int64, message tgbotapi.Message) {
    countdown := 10
    // Обновляем сообщение каждую секунду
    for countdown > 0 {
        // Формируем текст с обратным отсчетом
        msgText := fmt.Sprintf("%s\nСообщение удалится через %d секунд...", message.Text, countdown)

        // Создаем запрос на редактирование сообщения
        editMsg := tgbotapi.NewEditMessageText(chatID, message.MessageID, msgText)

        // Отправляем запрос на редактирование сообщения
        _, err := BotAPI.Send(editMsg)
        if err != nil {
            fmt.Println("Ошибка при редактировании сообщения:", err)
            return
        }

        // Уменьшаем таймер обратного отсчета
        countdown--

        // Ждем 1 секунду
        time.Sleep(1 * time.Second)
    }

    // Удаление сообщения после обратного отсчета
    messageWithDesc := tgbotapi.NewDeleteMessage(chatID, message.MessageID)
    _, delErr := BotAPI.Send(messageWithDesc)
    if delErr != nil {
        fmt.Println("Ошибка при удалении сообщения:", delErr)
    }
}

