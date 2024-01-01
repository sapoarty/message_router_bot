package database

import (
    "message_router_bot/utils"
    "log"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "strings"
)

func AddKeywords(keywordsList[] string, chatID int64, userID int) (error) {
    if _, ok := utils.UsersKeywordsChatsMap[userID]; !ok {
        utils.InitUsersKeywordsChatsMap(userID)
    }

    sqlStmt := `INSERT OR REPLACE INTO keywords (user_id, keyword, chat_id) values (?,?,?)`
    for _, keyword := range keywordsList {
        _, err := DB.Exec(sqlStmt, userID, keyword, chatID)
        if err != nil {
            log.Printf("SetKeywords err %d", err) 
            return err 
        }
        utils.UsersKeywordsChatsMap[userID][keyword] = chatID
    }
    log.Printf(
        "Keywords list [%s] added for chatID %d\n", 
        strings.Join(keywordsList, ","), 
        chatID,
    )
    return nil
}

func DeleteKeywords(keywordsList []string, chatID int64, userID int) ([]string, error) {
    sqlStmt := "DELETE FROM keywords WHERE user_id = ? AND keyword = ? AND chat_id = ?"
    log.Printf("DeleteKeywords %s", sqlStmt)
    
    var notFoundKeywords []string
    for _, keyword := range keywordsList {
        // Добавляем проверку наличия keyword в БД
        err := DB.QueryRow("SELECT 1 FROM keywords WHERE user_id = ? AND keyword = ?", userID, keyword).Scan(new(int))
        if err != nil {
            if err == sql.ErrNoRows {
                // Если keyword не найден, добавляем его в список не найденных
                notFoundKeywords = append(notFoundKeywords, keyword)
            } else {
                // Если произошла другая ошибка - возвращаем ее
                log.Printf("Error during checking keyword %s: %s", keyword, err)
                return nil, err
            }
        } else {
            _, err := DB.Exec(sqlStmt, userID, keyword, chatID)
            if err != nil {
                log.Printf("DeleteKeywords %s, err %s", sqlStmt, err)
                return nil, err
            }
            utils.DeleteKeywordFromLocalStore(chatID, keyword, userID)
        }
    }
    log.Printf(
        "Keywords list [%s] were tried to delete for chatID %d\n",
        strings.Join(keywordsList, ","),
        chatID,
    )
    return notFoundKeywords, nil
}


func DeleteChatData(chatID int64, userID int) (error) {
    sqlStmt := "DELETE FROM keywords WHERE user_id = ? AND chat_id = ?"
    log.Printf("DeleteKeywords %s", sqlStmt)
    _, err := DB.Exec(sqlStmt, userID, chatID)
    if err != nil {
        log.Printf("DeleteKeywords %s, err %s", sqlStmt, err)
        return err
    }
    utils.DeleteChatFromLocalStore(chatID, userID)
    log.Printf(
        "chatID [%s] data was deleted\n",
        chatID,
    )
    return nil
}