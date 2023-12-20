package database

import (
    "log"
    // "fmt"
    _ "github.com/mattn/go-sqlite3"
    "encoding/json"
    "strings"
)

func AddKeywords(keywordsList[] string, chatID int64, userID int) (error) {
    sqlStmt := `INSERT OR IGNORE INTO keywords (user_id, keyword, chat_id) values (?,?,?)`

    for _, keyword := range keywordsList {
        _, err := DB.Exec(sqlStmt, userID, keyword, chatID)
        if err != nil {
            log.Printf("SetKeywords err %d", err) 
            return err 
        }
        if _, ok := UsersKeywordsChatsMap[userID]; !ok {
            UsersKeywordsChatsMap[userID] = make(map[string]int64)
        }
        UsersKeywordsChatsMap[userID][keyword] = chatID
    }
    log.Printf(
        "Keywords list [%s] added for chatID %d\n", 
        strings.Join(keywordsList, ","), 
        chatID,
    )
    return nil
}

func DeleteKeywords(keywordsList[] string, chatID int64, userID int) (error) {
    sqlStmt := `Delete from keywords where user_id = ? keyword = ? and chat_id = ?`
    log.Printf("SetKeywords %s", sqlStmt) 
    for _, keyword := range keywordsList {
        _, err := DB.Exec(sqlStmt, userID, keyword, chatID)
        if err != nil {
            log.Printf("SetKeywords %s, err %s", sqlStmt, err) 
            return err 
        }
        deleteKeyword(chatID, keyword, userID)
    }
    log.Printf(
        "Keywords list [%s] delete for chatID %d\n", 
        strings.Join(keywordsList, ","), 
        chatID,
    )
    return nil
}


func PrintUsersKeywordsChatsMap() (error) {
    log.Println("UsersKeywordsChatsMap: ")
    kwJson, err := json.Marshal(UsersKeywordsChatsMap)
    if err != nil {
        log.Println(err)
        return err
    }

    log.Println(string(kwJson)) 
    return nil
}

func GetKeywordsForUserChatID(chatId int64, userID int) ([]string) {
    var keywordsList []string
    for curKeyword, curChatId := range UsersKeywordsChatsMap[userID] {
        if (curChatId == chatId) {
            keywordsList = append(keywordsList, curKeyword)
        }
    }
    return keywordsList
}

func deleteKeyword(chatID int64, keyword string, userID int) {
    if keywordsListForUserId, ok := UsersKeywordsChatsMap[userID]; ok {
        delete(keywordsListForUserId, keyword)
    }
}