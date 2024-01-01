package database

import (
    "message_router_bot/structures"
    "log"
    _ "github.com/mattn/go-sqlite3"
    "github.com/jmoiron/sqlx"
    "message_router_bot/utils"
)

var DB *sqlx.DB

func Init() {
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
    utils.PrintUsersKeywordsChatsMap(0)
}

// LoadKeywordsData загружает данные о ключевых словах из базы данных и помещает их в глобальную карту keywordChatMap.
func LoadKeywordsData() (error) {
    var keywords []structures.Keyword
    err := DB.Select(&keywords, "SELECT * FROM keywords")
    if err != nil {
        return err
    }

    for _, keyword := range keywords {
        if _, ok := utils.UsersKeywordsChatsMap[keyword.UserID]; !ok {
            utils.InitUsersKeywordsChatsMap(keyword.UserID)
        }
        utils.UsersKeywordsChatsMap[keyword.UserID][keyword.Keyword] = keyword.ChatID
    }
    return nil
}

