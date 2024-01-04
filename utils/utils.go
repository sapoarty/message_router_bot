package utils

import (
	"message_router_bot/config"
    "message_router_bot/structures"
    "message_router_bot/constants"
    "log"
    "strings"
    "regexp"
)


func InitUserData(userID int) {
    log.Printf("InitUserData for user %d", userID)
    if config.UserStates == nil {
        config.UserStates = make(map[int]structures.User)
        log.Println("Init UserStates with null")
    }
    config.UserStates[userID] = structures.User{ExpectingInput: false, UserCommand: "", Lang: "ru"}
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
    keywords := structures.GetKeywordsForUserChatID(chatID, userID)
    return (len(keywords) > 0 && keywords[0] == constants.DefaultGroup)
}

func SanitizeKeywords(keywordsList []string) []string {
    var sanitizedKeywords []string
    // Регулярное выражение, определяющее допустимые символы в ключевых словах
    re := regexp.MustCompile("^[a-zA-Z0-9_]+$")

    for _, keyword := range keywordsList {
        // Удаление начальных и конечных пробелов
        trimmedKeyword := strings.TrimSpace(keyword)
        // Дополнительная проверка на длину ключевого слова
        if len(trimmedKeyword) >= 1 && len(trimmedKeyword) <= 1024 {
            // Проверка ключевого слова на соответствие допустимым символам
            if re.MatchString(trimmedKeyword) {
                sanitizedKeywords = append(sanitizedKeywords, trimmedKeyword)
            }
        }
    }

    return sanitizedKeywords
}

func IsMessageTextContainKeyword(origMsgText string,keyword string) bool {
    return origMsgText != "" && strings.Contains(strings.ToLower(origMsgText), strings.ToLower(keyword))
}