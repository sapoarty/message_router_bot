package structures

import (
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyword struct {
	ID       int64      `db:"id" json:"id"`
	UserID   int        `db:"user_id" json:"user_id"`
	ChatID   int64      `db:"chat_id" json:"chat_id"`
	Keyword  string     `db:"keyword" json:"keyword"`
}

type User struct {
	ExpectingInput bool
	UserCommand    string
	Lang		   string
}

type BotMessagesHandler struct {
	Handler func(text string, chatID int64, userID int)
	NeedReply bool
}

type MediaGroup struct {
	Messages []*tgbotapi.Message // Сообщения в группе
	Timer    *time.Timer         // Таймер для отправки группы сообщений
}

type UserMediaGroups struct {
	Groups map[string]*MediaGroup // Карта групп медиафайлов по groupID
}