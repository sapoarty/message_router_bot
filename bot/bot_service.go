package bot

import (
    "message_router_bot/constants"
    "message_router_bot/messages"
    "message_router_bot/structures"
    "message_router_bot/config"
    "message_router_bot/utils"
    "message_router_bot/database"
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
    "fmt"
    "strings"
    "time"
)

var handlerMap map[string]structures.BotMessagesHandler

func ForwardMessage(messagesList[] *tgbotapi.Message) error {
    if (len(messagesList) == 0) {
        return nil
    }

    log.Printf("ForwardMessage cnt = %d", len(messagesList))

    var matchFound = false
    var usedKeywordsForChatsList = make(map[int64]string)
    var message = messagesList[0]
    var mediaGroup []interface{}
    var origMsgText string

    if (len(messagesList) > 1) {

        for _, msg := range messagesList {
            if msg.Photo != nil {
                biggestPhotoSize := msg.Photo[len(msg.Photo)-1]
                media := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(biggestPhotoSize.FileID))
                if msg.Caption != "" {
                    log.Println("msg.Caption: " + msg.Caption)
                    media.Caption = msg.Caption
                    origMsgText += msg.Caption
                }
                mediaGroup = append(mediaGroup, media)
                messageToDel := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
                _, delErr := BotAPI.Request(messageToDel)
                if delErr != nil {
                    fmt.Println("Ошибка при удалении сообщения:", delErr)
                }

            }
        }
    }

    userID := int(message.From.ID)
    botChatID := message.Chat.ID
    userLang := config.UserStates[userID].Lang

    if origMsgText == "" {
        origMsgText = strings.TrimSpace(message.Text)            
    }
    if origMsgText == "" {
        origMsgText = strings.TrimSpace(message.Caption)
    }
    urls := ExtractHiddenURLs(message)
    if (len(urls) > 0) {
        log.Println("ExtractHiddenURLs" + strings.Join(urls, ", "))
        origMsgText += messages.HidedUrlInMessage[userLang] + strings.Join(urls, ", ")
    }

    if utils.IsMessageTextContainKeyword(origMsgText, "http") {
        desc, err := utils.GetMetaDescription(utils.GetURLFrom(origMsgText))
        if err != nil {
            log.Printf("Err GetMetaDescription: %s\n", err)
            return err
        }
        origMsgText += messages.DescUrl[userLang] + desc
        log.Println("desc: ", desc)
    }

    // Итерируемся через наши ключевые слова и их идентификаторы чата
    for keyword, chatID := range structures.UsersKeywordsChatsMap[userID] {
        // Если в этот чат уже отправляли, не повторяем
        if _, ok := usedKeywordsForChatsList[chatID]; ok {
            continue
        }
        // Если описание или сообщение содержит ключевое слово (без учета регистра)
        if utils.IsMessageTextContainKeyword(origMsgText, keyword) {
            err, groupName := GetGroupNameByChatId(chatID)
            if err != nil {
                log.Printf("Err GetGroupNameByChatId: %s\n", err)
                return err
            }
            log.Println("Contains keyword:", keyword, "Group Name:", groupName)

            // Создаем запрос на пересылку сообщения
            var msgToResend tgbotapi.Chattable
            if len(messagesList) > 1 {
                log.Printf("Files in group %d", len(mediaGroup))
                msgToResend = tgbotapi.NewMediaGroup(chatID, mediaGroup)
            } else {
                msgToResend = tgbotapi.NewCopyMessage(chatID, botChatID, message.MessageID)                
            }
            _, err = BotAPI.Send(msgToResend)
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
            sendMessage(messages.DefaulGroupIsNotSet[userLang], botChatID)
            return nil
        }
        // Создаем запрос на пересылку сообщения в группу по-умолчанию
        var msgToResend tgbotapi.Chattable
        if len(messagesList) > 1 {
            log.Printf("Files in group %d", len(mediaGroup))
            msgToResend = tgbotapi.NewMediaGroup(chatID, mediaGroup)
            BotAPI.Send(msgToResend)
        } else {
            msgToResend = tgbotapi.NewCopyMessage(chatID, botChatID, message.MessageID)                
            _, err := BotAPI.Send(msgToResend)
            if err != nil {
                log.Println("Error during forwarding message to default group:", err)
                return err
            }
        }
    }

    if (matchFound) {
        matchFoundMessageSend(usedKeywordsForChatsList, origMsgText, botChatID, userID)
    } else {
        matchNotFoundMessageSend(origMsgText, botChatID, userID)
    }
    
    return nil
}

func matchFoundMessageSend(usedKeywordsForChatsList map[int64]string, origMsgText string, botChatID int64, userID int) (error) {
    userLang := config.UserStates[userID].Lang
    var chatNamesListStr string
    for chat, keyword := range usedKeywordsForChatsList {
        _, chatName := GetGroupNameByChatId(chat)
        chatNamesListStr += fmt.Sprintf("%s(%s), ", chatName, keyword)
    }
    chatNamesListStr = strings.TrimSuffix(chatNamesListStr, ", ")
    runes := []rune(origMsgText)
    if len(runes) > constants.ForwardMessageCutLen {
        origMsgText = string(runes[:constants.ForwardMessageCutLen]) + "..."
    }
    msgText := fmt.Sprintf(messages.MessageHasBeenForwardedToChatsUsingKeywords[userLang], origMsgText, chatNamesListStr)
    resultMsg, err := sendMessage(msgText, botChatID, "HTML")
    if (err != nil) {
        return err
    }
    go countdownAndDeleteMessage(botChatID, resultMsg)
    return nil
}

func matchNotFoundMessageSend(origMsgText string, botChatID int64, userID int) (error) {
    userLang := config.UserStates[userID].Lang
    defaultChatID := structures.GetUserChatIDForKeyword(constants.DefaultGroup, userID)
    _, defaultChatName := GetGroupNameByChatId(defaultChatID)
    chatID := structures.GetUserChatIDForKeyword(constants.DefaultGroup, userID)
    if chatID == 0 {
        sendMessage(messages.DefaulGroupIsNotSet[userLang], botChatID)
        return fmt.Errorf("Default Group is not set error")
    }
    msgText := fmt.Sprintf(messages.MessageHasBeenForwardedToDefaultGroup[userLang], origMsgText, defaultChatName)
    resultMsg, err := sendMessage(msgText, botChatID, "HTML")
    if (err != nil) {
        return err
    }
    log.Println("2 " + msgText)
    go countdownAndDeleteMessage(botChatID, resultMsg)
    return nil
}

func askToReply(message *tgbotapi.Message) {
    text := strings.TrimSpace(message.Text)
    chatID := message.Chat.ID
    userID := int(message.From.ID)
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

func GetGroupNameByChatId(chatId int64) (error, string) {
    chatConfig := tgbotapi.ChatInfoConfig{
        ChatConfig: tgbotapi.ChatConfig{
            ChatID: chatId,
        },
    }

    chat, err := BotAPI.GetChat(chatConfig)
    if err != nil {
        return err, ""
    }

    // Возвращаем title чата, который для группы будет названием группы.
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
        msgText := fmt.Sprintf("%s\n\n<b>Сообщение удалится через %d секунд...</b>", message.Text, countdown)

        // Создаем запрос на редактирование сообщения
        editMsg := tgbotapi.NewEditMessageText(chatID, message.MessageID, msgText)
        editMsg.ParseMode = tgbotapi.ModeHTML
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
    _, delErr := BotAPI.Request(messageWithDesc)
    if delErr != nil {
        fmt.Println("Ошибка при удалении сообщения:", delErr)
    }
}


func ExtractHiddenURLs(message *tgbotapi.Message) []string {
    urls := []string{}

    // Обработка сущностей сообщения
    for _, entity := range message.Entities {
        if entity.Type == "text_link" && entity.URL != "" {
            urls = append(urls, entity.URL)
        }
    }

    // Обработка сущностей подписи, если это сообщение с медиа
    if message.Caption != "" {
        for _, entity := range message.CaptionEntities {
            if entity.Type == "text_link" && entity.URL != "" {
                urls = append(urls, entity.URL)
            }
        }
    }

    return urls
}