package structures

import (
	"encoding/json"
	"log"
	"sync"
)

// Структура для хранения id и имен чатов и каналов
var UsersKeywordsChatsMap = make(map[int]map[string]int64)
var mu sync.RWMutex // Мьютекс для синхронизации доступа к UsersKeywordsChatsMap

func InitUsersKeywordsChatsMap(userID int) {
	mu.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer mu.Unlock()
	UsersKeywordsChatsMap[userID] = make(map[string]int64)
}

func DeleteKeywordFromLocalStore(chatID int64, keyword string, userID int) {
	mu.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer mu.Unlock()
	if keywordsListForUserId, ok := UsersKeywordsChatsMap[userID]; ok {
		delete(keywordsListForUserId, keyword)
	}
}

func DeleteChatFromLocalStore(chatID int64, userID int) {
	mu.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer mu.Unlock()
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
	mu.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer mu.Unlock()    
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

	mu.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer mu.Unlock()
	for curKeyword, curChatId := range UsersKeywordsChatsMap[userID] {
		if (curChatId == chatID) {
			keywordsList = append(keywordsList, curKeyword)
		}
	}
	return keywordsList
}

func GetUserChatIDForKeyword(keyword string, userID int) (int64) {
	mu.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer mu.Unlock()
	return UsersKeywordsChatsMap[userID][keyword] 
}

func GetKeywordsMapForUser(userID int) (map[string]int64) {
	mu.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer mu.Unlock()
	return UsersKeywordsChatsMap[userID]
}