package bot

import (
    "message_router_bot/constants"
    "message_router_bot/structures"
    "message_router_bot/config"
    "message_router_bot/utils"
    "message_router_bot/database"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "fmt"
    "strings"
)

var handlerMap map[string]structures.BotMessagesHandler

func ForwardMessage(message *tgbotapi.Message) error {
    var matchFound = false
    var usedChatsList = make(map[int64]bool)
    userID := message.From.ID
    botChatID := message.Chat.ID
    messageToDelete := tgbotapi.NewDeleteMessage(botChatID, message.MessageID)

    text := ""
    if message.Text != "" { // Проверяем, что сообщение содержит текст
        text = message.Text
        log.Println(text)
    }

    // Итерируемся через наши ключевые слова и их идентификаторы чата
    for keyword, chatID := range utils.UsersKeywordsChatsMap[userID] {
        if (usedChatsList[chatID] == true) {
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
                    msg := fmt.Sprintf(constants.ChatDeleted[config.UserStates[userID].Lang], groupName, keyword)
                    sendMessage(msg, botChatID)
                    database.DeleteChatData(chatID, userID)
                } else {
                    log.Println("Error during forwarding message:", err)
                    return err
                }
            } else {
                matchFound = true
                usedChatsList[chatID] = true
            }
        }
    }

    if !matchFound {
        // Если в сообщении не нашлось ключевого слова, отправляем в группу по-умолчанию
        chatID := utils.GetUserChatIDForKeyword(constants.DefaultGroup, userID)
        if chatID == 0 {
            sendErrorMessage := tgbotapi.NewMessage(botChatID, constants.DefaulGroupIsNotSet[config.UserStates[userID].Lang])
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
        constants.InputKeywordsListRequest[userLang], 
        chatID,
        tgbotapi.ReplyKeyboardRemove{
            RemoveKeyboard: true,
            Selective:      false,
        },
    )
}

func printKeywordsToGroup(chatID int64, userID int) {
    lang := config.UserStates[userID].Lang
    err, groupName := GetGroupNameByChatId(chatID)
    if err != nil {
        log.Printf("Err GetGroupNameByChatId: %s\n", err)
        return
    }

    log.Printf("printKeywordsToGroup %s", groupName)
    keywords := utils.GetKeywordsForUserChatID(chatID, userID)
    if err != nil {
        log.Printf("Err GetKeywordsByGroupName: %s\n", err)
        return
    }

    var msgText string
    if len(keywords) == 0 {
        msgText = constants.KeywordsListEmpty[lang]
    } else {
        msgText = fmt.Sprintf(constants.KeywordsListForGroup[lang], groupName, strings.Join(keywords, ", "))
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

func sendMessage(msgText string, chatID int64, additionalArgs ...interface{}) error {
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
    _, err := BotAPI.Send(message)
    if err != nil {
        log.Printf("Ошибка при отправке сообщения в чат ID %d: %v", chatID, err)
        return err
    }

    return nil
}

func setGroupDesc(chatID int64, userID int) {
    lang := config.UserStates[userID].Lang
    log.Printf("setGroupDesc chatID %d, userID %d", chatID, userID)

    var descText string
    if (utils.IsGroupDefault(chatID, userID)) {
        descText = fmt.Sprintf(constants.DefaultGroupDesc[lang])
    } else {
        keywords := utils.GetKeywordsForUserChatID(chatID, userID)
        if (len(keywords) > 0) {
            descText = fmt.Sprintf(
                constants.GroupDesc[lang],
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


