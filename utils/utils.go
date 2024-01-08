package utils

import (
    "message_router_bot/config"
    "message_router_bot/structures"
    "message_router_bot/constants"
    "log"
    "strings"
    "regexp"
    "net/http"
    "golang.org/x/net/html"
    "golang.org/x/net/html/charset"
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

    for _, keyword := range keywordsList {
        // Удаление начальных и конечных пробелов
        trimmedKeyword := strings.TrimSpace(keyword)
        // Дополнительная проверка на длину ключевого слова
        if len(trimmedKeyword) >= 1 && len(trimmedKeyword) <= 1024 {
            sanitizedKeywords = append(sanitizedKeywords, trimmedKeyword)
        }
    }

    return sanitizedKeywords
}

func IsMessageTextContainKeyword(origMsgText string,keyword string) bool {
    return origMsgText != "" && strings.Contains(strings.ToLower(origMsgText), strings.ToLower(keyword))
}

// GetURLFrom извлекает первый найденный URL из текстового сообщения.
func GetURLFrom(text string) string {
    // Регулярное выражение для поиска URL начинается с http:// или https://
    urlRegex := regexp.MustCompile(`https?://\S+`)
        // Найти все совпадения в строке
        matches := urlRegex.FindStringSubmatch(text)
        if len(matches) == 0 {
            return "" // Если URL не найден, возвращаем пустую строку
        }
        // Возвращает первое найденное совпадение URL
        return matches[0]
    }

    func GetMetaDescription(url string) (string, error) {
        resp, err := http.Get(url)
        if err != nil {
            return "", err
        }
        defer resp.Body.Close()

        // Принудительное преобразование в UTF-8, если сервер не отправляет нужные заголовки
        utf8Body, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
        if err != nil {
            return "", err
        }

        z := html.NewTokenizer(utf8Body)

        var title, description string
        var titleFound, metaFound bool

        for {
            tt := z.Next()

            switch tt {
            case html.ErrorToken:
                // Конец документа или ошибка, если ничего не найдено, возвращаем ошибку или недостающие данные
                if !titleFound || !metaFound {
                    return title + "\n" + description, nil
                }
            case html.StartTagToken, html.SelfClosingTagToken:
                t := z.Token()
                if t.Data == "title" && !titleFound {
                    tt = z.Next()
                    if tt == html.TextToken {
                        title = string(z.Text())
                        titleFound = true
                        if metaFound {
                            return title + "\n" + description, nil
                        }
                    }
                } else if t.Data == "meta" {
                    desc := ""
                    isDescription := false
                    for _, a := range t.Attr {
                        if a.Key == "name" && a.Val == "description" {
                            isDescription = true
                        } else if a.Key == "content" {
                            desc = a.Val
                        }

                        if isDescription && desc != "" {
                            description = desc
                            metaFound = true
                            if titleFound {
                                return title + "\n" + description, nil
                            }
                            break
                        }
                    }
                }
            }
        }

        return title + "\n" + description, nil
    }