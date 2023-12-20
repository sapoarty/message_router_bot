package database

import (
    "log"
    _ "github.com/mattn/go-sqlite3"
    "github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

type Keyword struct {
    ID       int64      `db:"id" json:"id"`
    UserID   int        `db:"user_id" json:"user_id"`
    ChatID   int64      `db:"chat_id" json:"chat_id"`
    Keyword  string     `db:"keyword" json:"keyword"`
}

// Структура для хранения id и имен чатов и каналов
var UsersKeywordsChatsMap = make(map[int]map[string]int64)

func InitDb() {
    var err error
    DB, err = sqlx.Connect("sqlite3", "./routing_bot_data.db")
    if err != nil {
        panic(err)
    }

    keywordStmt := `CREATE TABLE IF NOT EXISTS keywords (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INT NOT NULL,
        keyword TEXT NOT NULL,
        chat_id INT NOT NULL,
        UNIQUE (user_id, keyword, chat_id)
    );`
    _, err = DB.Exec(keywordStmt)
    if err != nil {
        log.Printf("%q: %s\n", err, keywordStmt)
        panic(err)
    }
    LoadKeywordsData()
    PrintUsersKeywordsChatsMap()
}

// LoadKeywordsData загружает данные о ключевых словах из базы данных и помещает их в глобальную карту keywordChatMap.
func LoadKeywordsData() (error) {
    var keywords []Keyword
    err := DB.Select(&keywords, "SELECT * FROM keywords")
    if err != nil {
        return err
    }

    for _, keyword := range keywords {
        if _, ok := UsersKeywordsChatsMap[keyword.UserID]; !ok {
            UsersKeywordsChatsMap[keyword.UserID] = make(map[string]int64)
        }
        UsersKeywordsChatsMap[keyword.UserID][keyword.Keyword] = keyword.ChatID
    }
    return nil
}

