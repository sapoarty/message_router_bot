package bot

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "message_router_bot/constants"
    "message_router_bot/database"
    "message_router_bot/config"
    "message_router_bot/utils"
    "fmt"
    "log"
    "strings"
    "os"
)

func HandleUpdates(updates tgbotapi.UpdatesChannel) {
    // Этот цикл просматривает все обновления, которые приходят от телеграмма
    for update := range updates {
        // Обработка CallbackQuery от inline кнопок
        if update.CallbackQuery != nil {
            log.Printf("CallbackQuery received: %s from user id %d", update.CallbackQuery.Data, update.CallbackQuery.From.ID)
            // Получаем данные от CallbackQuery
            cq := update.CallbackQuery
            callBackData := cq.Data
            callbackQueryID := cq.ID
            userID := cq.From.ID
            chatID := cq.Message.Chat.ID
            // Если нет зарегистрированного состояния для этого пользователя, создаем новое
            if _, ok := config.UserStates[userID]; !ok {
                utils.InitUserData(userID)
            }

            // Определяем, на какую inline кнопку было совершено нажатие используя callBackData
            handleCategorySelection(callBackData, chatID, userID)

            // Обратный вызов API для уведомления Telegram, что CallbackQuery был получен и обработан
            callbackConfig := tgbotapi.NewCallback(callbackQueryID, "")
            if _, err := BotAPI.AnswerCallbackQuery(callbackConfig); err != nil {
                log.Panic(err)
            }
            // Проверяем, есть ли сообщение в обновлении
        } else if update.Message != nil {
            log.Printf("Message %s from user id %d", update.Message.Text, update.Message.From.ID)
            message := update.Message
            chatID := message.Chat.ID
            userID := message.From.ID
            text := strings.TrimSpace(message.Text)
            chatType := message.Chat.Type
            // Если нет зарегистрированного состояния для этого пользователя, создаем новое
            if _, ok := config.UserStates[userID]; !ok {
                utils.InitUserData(userID)
                handlerMap = GetHandlerMapForUser(userID)
            }

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
                for _, newUser := range *update.Message.NewChatMembers {
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
        constants.AwaitingCommand[config.UserStates[userID].Lang], 
        os.Getenv("BotLink"),
    )
    sendMessage(msg, chatID, GetMenu(userID), "Markdown")
}