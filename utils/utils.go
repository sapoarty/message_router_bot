package utils

import (
	"message_router_bot/config"
    "message_router_bot/structures"
    "message_router_bot/constants"
    "encoding/json"
    "log"
    "strings"
)

// Структура для хранения id и имен чатов и каналов
var UsersKeywordsChatsMap = make(map[int]map[string]int64)

func InitUserData(userID int) {
    log.Printf("InitUserData for user %d", userID)
    if config.UserStates == nil {
        config.UserStates = make(map[int]structures.User)
        log.Println("Init UserStates with null")
    }
    config.UserStates[userID] = structures.User{ExpectingInput: false, UserCommand: "", Lang: "ru"}
}

func InitUsersKeywordsChatsMap(userID int) {
    UsersKeywordsChatsMap[userID] = make(map[string]int64)
}

func SetUserData(userID int, expectingInput *bool, userCommand *string, lang *string) {
    userState := config.UserStates[userID]
    
    if expectingInput != nil {
        userState.ExpectingInput = *expectingInput
    }
    
    if userCommand != nil {
        userState.UserCommand = *userCommand
    }
    
    if lang != nil {
        userState.Lang = *lang
    }
    
    config.UserStates[userID] = userState
}

func DeleteKeywordFromLocalStore(chatID int64, keyword string, userID int) {
    if keywordsListForUserId, ok := UsersKeywordsChatsMap[userID]; ok {
        delete(keywordsListForUserId, keyword)
    }
}

func DeleteChatFromLocalStore(chatID int64, userID int) {
    if keywordsListForUserId, ok := UsersKeywordsChatsMap[userID]; ok {
        for curKeyword, curChatId := range UsersKeywordsChatsMap[userID] {
            if (curChatId == chatID) {
                delete(keywordsListForUserId, curKeyword)
            }
        }
    }
}

func PrintUsersKeywordsChatsMap(userID int) (error) {
    var kwJson []byte
    var err error

    log.Println("UsersKeywordsChatsMap: ")
    if (userID == 0) {
        kwJson, err = json.Marshal(UsersKeywordsChatsMap)
        if err != nil {
            log.Println(err)
            return err
        }
    } else {
        kwJson, err = json.Marshal(UsersKeywordsChatsMap[userID])
        if err != nil {
            log.Println(err)
            return err
        }
    }

    log.Println(string(kwJson)) 
    return nil
}

func GetKeywordsForUserChatID(chatID int64, userID int) ([]string) {
    var keywordsList []string
    for curKeyword, curChatId := range UsersKeywordsChatsMap[userID] {
        if (curChatId == chatID) {
            keywordsList = append(keywordsList, curKeyword)
        }
    }
    return keywordsList
}

func GetUserChatIDForKeyword(keyword string, userID int) (int64) {
    return UsersKeywordsChatsMap[userID][keyword] 
}

func GetKeywordsMapForUser(userID int) (map[string]int64) {
    return UsersKeywordsChatsMap[userID]
}

func TrimSpacesFromStringList(keywordsList []string) ([]string){
    trimmedKeywordsList := make([]string, 0, len(keywordsList))
    for _, keyword := range keywordsList {
        if trimmedKeyword := strings.TrimSpace(keyword); trimmedKeyword != "" {
            trimmedKeywordsList = append(trimmedKeywordsList, trimmedKeyword)
        }
    }
    return trimmedKeywordsList
}

func ToLowercaseSlice(slice []string) []string {
    lowercaseSlice := make([]string, len(slice))
    for i, s := range slice {
        lowercaseSlice[i] = strings.ToLower(s)
    }
    return lowercaseSlice
}

func IsGroupDefault(chatID int64, userID int) bool {
    keywords := GetKeywordsForUserChatID(chatID, userID)
    return (len(keywords) > 0 && keywords[0] == constants.DefaultGroup)
}