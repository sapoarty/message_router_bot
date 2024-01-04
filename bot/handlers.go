package bot

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "message_router_bot/constants"
    "message_router_bot/messages"
    "message_router_bot/structures"
    "message_router_bot/database"
    "message_router_bot/config"
    "message_router_bot/utils"
    "log"
    "fmt"
    "strings"
    "os"
)

func GetHandlerMapForUser (userID int) (map[string]structures.BotMessagesHandler) {
    userLang := config.UserStates[userID].Lang
    return map[string]structures.BotMessagesHandler{
        messages.CommandHelp[userLang]: {Handler: handleStart, NeedReply: false},
        messages.CommandStart[userLang]: {Handler: handleStart, NeedReply: false},
        messages.CommandAddByCategory[userLang]: {Handler: handleShowCategoriesInlineKeybord, NeedReply: false},
        messages.CommandAddKeywords[userLang]: {Handler: handleAddKeyWordsForGroup, NeedReply: true},
        messages.CommandDeleteKeywords[userLang]: {Handler: handleDeleteKeyWordsForGroup, NeedReply: true},
        messages.CommandSetDefaultGroup[userLang]: {Handler: handleSetDefaultGroup, NeedReply: false},
        messages.CommandPrintAllKeywords[userLang]: {Handler: handlePrintAllKeywords, NeedReply: false},
        messages.CommandChangeLang[userLang]: {Handler: handleChangeLang, NeedReply: false},
    }   
}

func handleStart(text string, chatID int64, userID int) {    
    log.Printf("handleStart text [%s], chatID %d, userID %d", text, chatID, userID)
    msg := fmt.Sprintf(
        messages.GreetingsMessage[config.UserStates[userID].Lang] + messages.AwaitingCommand[config.UserStates[userID].Lang], 
        os.Getenv("BotLink"),
    )
    sendMessage(msg, chatID,GetMenu(userID),"Markdown")
}


func handleShowCategoriesInlineKeybord(text string, chatID int64, userID int) {
    var rows [][]tgbotapi.InlineKeyboardButton
    lang := config.UserStates[userID].Lang

    msg := tgbotapi.NewMessage(chatID, messages.ChooseCategory[lang])
    categoriesMap := constants.GetDefaultCategories()
    for categoryName, _ := range categoriesMap {
        button := tgbotapi.NewInlineKeyboardButtonData(categoryName, categoryName)
        row := tgbotapi.NewInlineKeyboardRow(button)
        rows = append(rows, row)
    }

    // Создаем Inline клавиатуру из rows
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

    // Send the message with the inline keyboard
    if _, err := BotAPI.Send(msg); err != nil {
        log.Panic(err)
    }
}

func handleAddKeyWordsForGroup(text string, chatID int64, userID int) {
    log.Printf("handleAddKeyWordsForGroup text [%s], chatID %d, userID %d", text, chatID, userID)
    userLang := config.UserStates[userID].Lang
    trimmedKeywordsList := utils.TrimSpacesFromStringList(strings.Split(text, ","))
    trimmedLowercaseKeywordsList := utils.ToLowercaseSlice(trimmedKeywordsList)

    if len(trimmedLowercaseKeywordsList) == 0 {
        sendMessage(messages.KeywordsListEmpty[userLang], chatID)
        return
    }

    err := database.AddKeywords(trimmedLowercaseKeywordsList, chatID, userID)
    if err != nil {
        log.Println("Err AddKeywords: ", err)
        return 
    }
    setGroupDesc(chatID, userID)
    keywords := structures.GetKeywordsForUserChatID(chatID, userID)
    msg := fmt.Sprintf(messages.KeywordsListIsChanged[userLang], strings.Join(keywords, ", "))
    sendMessage(msg, chatID)
}

func handleDeleteKeyWordsForGroup(text string, chatID int64, userID int) {
    log.Printf("handleDeleteKeyWordsForGroup text [%s], chatID %d, userID %d", text, chatID, userID)
    userLang := config.UserStates[userID].Lang
    if (utils.IsGroupDefault(chatID, userID)) {
        sendMessage(messages.GroupIsAlreadyDefaultGroup[userLang], chatID)
        return
    }
    trimmedKeywordsList := utils.TrimSpacesFromStringList(strings.Split(text, ","))
    trimmedLowercaseKeywordsList := utils.ToLowercaseSlice(trimmedKeywordsList)

    if len(trimmedLowercaseKeywordsList) == 0 {
        sendMessage(messages.KeywordsListEmpty[userLang], chatID)
        return
    }

    result, err := database.DeleteKeywords(trimmedLowercaseKeywordsList, chatID, userID)
    if err != nil {
        log.Println("Err DeleteKeywords: ", err)
        return 
    }

    if len(result) > 0 {
        sendMessage(
            fmt.Sprintf(messages.KeywordsNotFound[userLang], strings.Join(result, ",")), 
            chatID,
        )
    } else {
        sendMessage(
            fmt.Sprintf(messages.AllKeywordsFoundAndDeleted[userLang]), 
            chatID,
        )
    }

    setGroupDesc(chatID, userID)
    keywords := structures.GetKeywordsForUserChatID(chatID, userID)
    msg := fmt.Sprintf(messages.KeywordsListIsChanged[userLang], strings.Join(keywords, ", "))
    sendMessage(msg, chatID)
}

func handleSetDefaultGroup(text string, chatID int64, userID int) {
    log.Printf("handleSetDefaultGroup text [%s], chatID %d, userID %d", text, chatID, userID)
    userLang := config.UserStates[userID].Lang
    keywords := structures.GetKeywordsForUserChatID(chatID, userID)
    if (utils.IsGroupDefault(chatID, userID)) {
        sendMessage(
            messages.GroupIsAlreadyDefaultGroup[userLang], 
            chatID,
            tgbotapi.ReplyKeyboardRemove{
                RemoveKeyboard: true,
                Selective:      false,
            },
        )
        return
    } else if (len(keywords) > 0) {
        msg := fmt.Sprintf(messages.GroupHasKeywordsAndCanNotBeDefault[userLang], strings.Join(keywords, ", "))
        sendMessage(msg, chatID)
        return
    }

    oldDefaultChatId := structures.GetUserChatIDForKeyword(constants.DefaultGroup, userID)
    if (oldDefaultChatId != 0) {
        err, OldGroupName := GetGroupNameByChatId(oldDefaultChatId)
        if err != nil {
            log.Printf("Err GetGroupNameByChatId: %s\n", err)
            return
        }
        _, err = database.DeleteKeywords([]string {constants.DefaultGroup}, chatID, userID)
        if err != nil {
            log.Println("Err DeleteKeywords: ", err)
            return 
        }
        msgText := fmt.Sprintf(messages.DefaulGroupResetSuccessfully[userLang], OldGroupName)
        sendMessage(msgText, chatID)
        setGroupDesc(oldDefaultChatId, userID)
    }
    err := database.AddKeywords([]string{constants.DefaultGroup}, chatID, userID)
    if err != nil {
        log.Printf("Err AddKeywords: %s\n", err)
        return
    }
    log.Printf("Default group set fot chatID %d userID %d", chatID, userID)

    err, groupName := GetGroupNameByChatId(chatID)
    if err != nil {
        log.Printf("Err GetGroupNameByChatId: %s\n", err)
        return
    }

    msgText := fmt.Sprintf(messages.DefaulGroupSetSuccessfully[userLang], groupName)
    sendMessage(
        msgText, 
        chatID, 
        GetShortMenu(userID),
    )
    setGroupDesc(chatID, userID)
}

func handlePrintAllKeywords(text string, chatID int64, userID int) {
    log.Printf("handlePrintAllKeywords text [%s], chatID %d, userID %d", text, chatID, userID)
    userLang := config.UserStates[userID].Lang

    keywordsMap := structures.GetKeywordsMapForUser(userID)
    var msg strings.Builder
    if (len(keywordsMap) == 0) {
        sendMessage(messages.KeywordsAreNotSetYet[userLang], chatID)
    }

    // Организуем структуру, чтобы сгруппировать ключевые слова по chatID
    chatKeywords := make(map[int64][]string)
    for keyword, chatId := range keywordsMap {
        if (keyword == constants.DefaultGroup) {
            keyword = messages.DefaulGroupAlias[userLang]
        }
        chatKeywords[chatId] = append(chatKeywords[chatId], keyword)
    }

    // Теперь итерируемся через chatKeywords для создания итогового сообщения
    for chatId, keywords := range chatKeywords {
        err, chatName := GetGroupNameByChatId(chatId)
        if err != nil {
            chatName = "Unknown"
        }
        // Добавляем имя чата в сообщение
        msg.WriteString("<b>" + chatName + "</b>" + ": ")
        // Добавляем ключевые слова, разделенные запятой
        for i, keyword := range keywords {
            if i > 0 {
                msg.WriteString(", ")
            }
            msg.WriteString(keyword)
        }
        msg.WriteString("\n\n") // Завершаем список ключевых слов для текущего чата переносом строки
    }

    // Отправляем сформированное сообщение
    sendMessage(msg.String(), chatID, "HTML")
}


func handleChangeLang(text string, chatID int64, userID int) {
    var userLang string
    if (config.UserStates[userID].Lang == "ru") {
        userLang = "en"
    } else {
        userLang = "ru"
    }
    utils.SetUserData(userID, nil, nil, &userLang)
    handlerMap = GetHandlerMapForUser(userID)

    sendMessage(messages.LangIsChanged[config.UserStates[userID].Lang], chatID, GetMenu(userID))
}