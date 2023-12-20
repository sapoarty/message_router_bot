package bot

import (
    "message_router_bot/database"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "strings"
)

func GetChannelNameByChatId(chatID int64) (error, string) {
    chat, err := BotAPI.GetChat(tgbotapi.ChatConfig{ChatID: chatID})
    if err != nil {
        log.Println("GetChannelNameByChatId err %s", err)
        return err, ""
    }
    return nil, chat.Title
}

func ForwardMessage(message *tgbotapi.Message) error {
    // chatID := update.Message.Chat.ID
    text := message.Text
    userID := message.From.ID

    // Итерируемся через наши ключевые слова и их идентификаторы чата
    log.Println(text)
    for keyword, chatID := range database.UsersKeywordsChatsMap[userID] {
        // Если сообщение содержит ключевое слово (без учета регистра)
        if strings.Contains(strings.ToLower(text), strings.ToLower(keyword)) {
            err, channelName := GetChannelNameByChatId(chatID)
            if err != nil {
                log.Printf("Err GetChannelNameByChatId: %s\n", err)
                return nil
            }
            log.Println(text, keyword, channelName)
            // Создаем новый объект сообщения
            msg := tgbotapi.NewMessage(chatID, text)

            // Пересылаем сообщение
            BotAPI.Send(msg)

            messageToDelete := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
            _, err = BotAPI.DeleteMessage(messageToDelete)
            if err != nil {
                // обработка ошибки
            }
        } 
    }
    return nil
}