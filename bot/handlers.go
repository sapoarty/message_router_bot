package bot

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "message_router_bot/constants"
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
        constants.CommandHelp[userLang]: {Handler: handleStart, NeedReply: false},
        constants.CommandStart[userLang]: {Handler: handleStart, NeedReply: false},
        constants.CommandAddByCategory[userLang]: {Handler: handleShowCategoriesInlineKeybord, NeedReply: false},
        constants.CommandAddKeywords[userLang]: {Handler: handleAddKeyWordsForGroup, NeedReply: true},
        constants.CommandDeleteKeywords[userLang]: {Handler: handleDeleteKeyWordsForGroup, NeedReply: true},
        constants.CommandSetDefaultGroup[userLang]: {Handler: handleSetDefaultGroup, NeedReply: false},
        constants.CommandPrintAllKeywords[userLang]: {Handler: handlePrintAllKeywords, NeedReply: false},
        constants.CommandChangeLang[userLang]: {Handler: handleChangeLang, NeedReply: false},
    }   
}

func handleStart(text string, chatID int64, userID int) {    
    log.Printf("handleStart text [%s], chatID %d, userID %d", text, chatID, userID)
    msg := fmt.Sprintf(
        constants.GreetingsMessage[config.UserStates[userID].Lang] + constants.AwaitingCommand[config.UserStates[userID].Lang], 
        os.Getenv("BotLink"),
    )
    sendMessage(msg, chatID,GetMenu(userID),"Markdown")
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

    sendMessage(constants.LangIsChanged[config.UserStates[userID].Lang], chatID, GetMenu(userID))
}

func handleAddKeyWordsForGroup(text string, chatID int64, userID int) {
    log.Printf("handleAddKeyWordsForGroup text [%s], chatID %d, userID %d", text, chatID, userID)
    userLang := config.UserStates[userID].Lang
    trimmedKeywordsList := utils.TrimSpacesFromStringList(strings.Split(text, ","))
    trimmedLowercaseKeywordsList := utils.ToLowercaseSlice(trimmedKeywordsList)

    if len(trimmedLowercaseKeywordsList) == 0 {
        sendMessage(constants.KeywordsListEmpty[userLang], chatID)
        return
    }

    err := database.AddKeywords(trimmedLowercaseKeywordsList, chatID, userID)
    if err != nil {
        log.Println("Err AddKeywords: ", err)
        return 
    }
    setGroupDesc(chatID, userID)
    keywords := structures.GetKeywordsForUserChatID(chatID, userID)
    msg := fmt.Sprintf(constants.KeywordsListIsChanged[userLang], strings.Join(keywords, ", "))
    sendMessage(msg, chatID)
}

func handleDeleteKeyWordsForGroup(text string, chatID int64, userID int) {
    log.Printf("handleDeleteKeyWordsForGroup text [%s], chatID %d, userID %d", text, chatID, userID)
    userLang := config.UserStates[userID].Lang
    if (utils.IsGroupDefault(chatID, userID)) {
        sendMessage(constants.GroupIsAlreadyDefaultGroup[userLang], chatID)
        return
    }
    trimmedKeywordsList := utils.TrimSpacesFromStringList(strings.Split(text, ","))
    trimmedLowercaseKeywordsList := utils.ToLowercaseSlice(trimmedKeywordsList)

    if len(trimmedLowercaseKeywordsList) == 0 {
        sendMessage(constants.KeywordsListEmpty[userLang], chatID)
        return
    }

    result, err := database.DeleteKeywords(trimmedLowercaseKeywordsList, chatID, userID)
    if err != nil {
        log.Println("Err DeleteKeywords: ", err)
        return 
    }

    if len(result) > 0 {
        sendMessage(
            fmt.Sprintf(constants.KeywordsNotFound[userLang], strings.Join(result, ",")), 
            chatID,
        )
    } else {
        sendMessage(
            fmt.Sprintf(constants.AllKeywordsFoundAndDeleted[userLang]), 
            chatID,
        )
    }

    setGroupDesc(chatID, userID)
    keywords := structures.GetKeywordsForUserChatID(chatID, userID)
    msg := fmt.Sprintf(constants.KeywordsListIsChanged[userLang], strings.Join(keywords, ", "))
    sendMessage(msg, chatID)
}

func handleSetDefaultGroup(text string, chatID int64, userID int) {
    log.Printf("handleSetDefaultGroup text [%s], chatID %d, userID %d", text, chatID, userID)
    userLang := config.UserStates[userID].Lang
    keywords := structures.GetKeywordsForUserChatID(chatID, userID)
    if (utils.IsGroupDefault(chatID, userID)) {
        sendMessage(
            constants.GroupIsAlreadyDefaultGroup[userLang], 
            chatID,
            tgbotapi.ReplyKeyboardRemove{
                RemoveKeyboard: true,
                Selective:      false,
            },
        )
        return
    } else if (len(keywords) > 0) {
        msg := fmt.Sprintf(constants.GroupHasKeywordsAndCanNotBeDefault[userLang], strings.Join(keywords, ", "))
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
        msgText := fmt.Sprintf(constants.DefaulGroupResetSuccessfully[userLang], OldGroupName)
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

    msgText := fmt.Sprintf(constants.DefaulGroupSetSuccessfully[userLang], groupName)
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
        sendMessage(constants.KeywordsAreNotSetYet[userLang], chatID)
    }

    // Организуем структуру, чтобы сгруппировать ключевые слова по chatID
    chatKeywords := make(map[int64][]string)
    for keyword, chatId := range keywordsMap {
        if (keyword == constants.DefaultGroup) {
            keyword = constants.DefaulGroupAlias[userLang]
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

func handleBotMessage(message *tgbotapi.Message) {
    text := strings.TrimSpace(message.Text)
    chatID := message.Chat.ID
    userID := message.From.ID
    userLang := config.UserStates[userID].Lang

    if (text == constants.CommandStart[userLang] || text == constants.CommandHelp[userLang]) {
        // Если приветственное сообщение или нужна помощь, выводим текст и короткое меню
        msg := fmt.Sprintf(
            constants.GreetingsMessage[config.UserStates[userID].Lang] + constants.AwaitingCommand[config.UserStates[userID].Lang], 
            os.Getenv("BotLink"),
        )
        sendMessage(msg, chatID, GetShortMenu(userID), "Markdown")
    } else if (text == constants.CommandChangeLang[userLang]) {
        // Смена языка
        var userLang string
        if (config.UserStates[userID].Lang == "ru") {
            userLang = "en"
        } else {
            userLang = "ru"
        }
        utils.SetUserData(userID, nil, nil, &userLang)
        handlerMap = GetHandlerMapForUser(userID)

        sendMessage(constants.LangIsChanged[config.UserStates[userID].Lang], chatID, GetShortMenu(userID))
    } else if (text == constants.CommandPrintAllKeywords[userLang]) {
        handlePrintAllKeywords(text, chatID, userID)
    } else {
        // Иначе просто пересылаем сообщение, куда нужно
        log.Println("ForwardMessage")
        ForwardMessage(message)
    }
}

func handleDefaultGroupMessage(message *tgbotapi.Message) {
    text := strings.TrimSpace(message.Text)
    chatID := message.Chat.ID
    userID := message.From.ID
    userLang := config.UserStates[userID].Lang

    if (text == constants.CommandStart[userLang] || text == constants.CommandHelp[userLang]) {
        // Если приветственное сообщение или нужна помощь, выводим текст и короткое меню
        msg := fmt.Sprintf(
            constants.GreetingsMessage[userLang] +
            constants.GroupIsAlreadyDefaultGroup[userLang] + "\n" +
            constants.AwaitingCommand[userLang], 
            os.Getenv("BotLink"),
        )
        sendMessage(msg, chatID, GetShortMenu(userID), "Markdown")
    } else if (text == constants.CommandPrintAllKeywords[userLang]) {
        handlePrintAllKeywords(text, chatID, userID)
    } else if (text == constants.CommandChangeLang[userLang]) {
        // Смена языка
        var userLang string
        if (config.UserStates[userID].Lang == "ru") {
            userLang = "en"
        } else {
            userLang = "ru"
        }
        utils.SetUserData(userID, nil, nil, &userLang)
        handlerMap = GetHandlerMapForUser(userID)

        sendMessage(constants.LangIsChanged[config.UserStates[userID].Lang], chatID, GetShortMenu(userID))
    }
}

func handleShowCategoriesInlineKeybord(text string, chatID int64, userID int) {
    var rows [][]tgbotapi.InlineKeyboardButton
    lang := config.UserStates[userID].Lang

    msg := tgbotapi.NewMessage(chatID, constants.ChooseCategory[lang])
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

func handleCategorySelection(category string, chatID int64, userID int) {
    // Получаем карту категорий.
    categoriesMap := constants.GetDefaultCategories()

    // Проверяем, есть ли подкатегории в выбранной категории.
    words, exists := categoriesMap[category]
    if !exists {
        log.Printf("Категория %s не найдена", category)
        return
    }

    categoryWordsString := strings.Join(words, ",")
    handleAddKeyWordsForGroup(categoryWordsString, chatID, userID)
}