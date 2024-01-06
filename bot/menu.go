package bot

import (
    "message_router_bot/messages"
    "message_router_bot/config"
    "message_router_bot/constants"
    "fmt"
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMenu(userID int) tgbotapi.ReplyKeyboardMarkup {
    userLang := config.UserStates[userID].Lang
    commands := []string{
        messages.CommandHelp[userLang],
        messages.CommandAddByCategory[userLang],
        messages.CommandAddKeywords[userLang],
        messages.CommandDeleteKeywords[userLang],
        messages.CommandSetDefaultGroup[userLang],
        messages.CommandChangeLang[userLang],
    }
    var keyboardButtons [][]tgbotapi.KeyboardButton
    for _, command := range commands {
        row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(command))
        keyboardButtons = append(keyboardButtons, row)
    }

    keyboard := tgbotapi.NewReplyKeyboard(keyboardButtons...)
    return keyboard
}

func GetShortMenu(userID int) tgbotapi.ReplyKeyboardMarkup {
    userLang := config.UserStates[userID].Lang
    commands := []string{
        messages.CommandHelp[userLang],
        messages.CommandPrintAllKeywords[userLang],
        messages.CommandChangeLang[userLang],
    }

    var keyboardButtons [][]tgbotapi.KeyboardButton
    for _, command := range commands {
        row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(command))
        keyboardButtons = append(keyboardButtons, row)
    }

    keyboard := tgbotapi.NewReplyKeyboard(keyboardButtons...)
    return keyboard
}

// GetDefaultCategories выбирает локализованные названия категорий исходя из заданного языка.
func GetDefaultCategories(userID int) (map[string][]string, error) {
    userLang := config.UserStates[userID].Lang
    labels, ok := constants.CategoryLabels[userLang]
    if !ok {
        return nil, fmt.Errorf("Languange is not supported: %s", userLang)
    }

    localizedCategories := make(map[string][]string)
    for key, words := range constants.UniversalKeywords {
        label, ok := labels[key]
        if !ok {
            return nil, fmt.Errorf("No category localisation: %s", key)
        }
        localizedCategories[label] = words
    }

    return localizedCategories, nil
}