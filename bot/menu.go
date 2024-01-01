package bot

import (
    "message_router_bot/constants"
    "message_router_bot/config"
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetMenu(userID int) tgbotapi.ReplyKeyboardMarkup {
    userLang := config.UserStates[userID].Lang
    commands := []string{
        constants.CommandHelp[userLang],
        constants.CommandAddByCategory[userLang],
        constants.CommandAddKeywords[userLang],
        constants.CommandDeleteKeywords[userLang],
        constants.CommandSetDefaultGroup[userLang],
        constants.CommandChangeLang[userLang],
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
        constants.CommandHelp[userLang],
        constants.CommandPrintAllKeywords[userLang],
        constants.CommandChangeLang[userLang],
    }

    var keyboardButtons [][]tgbotapi.KeyboardButton
    for _, command := range commands {
        row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(command))
        keyboardButtons = append(keyboardButtons, row)
    }

    keyboard := tgbotapi.NewReplyKeyboard(keyboardButtons...)
    return keyboard
}
