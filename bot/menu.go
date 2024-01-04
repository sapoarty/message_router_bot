package bot

import (
    "message_router_bot/messages"
    "message_router_bot/config"
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
