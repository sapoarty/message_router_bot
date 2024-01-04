package structures

import (
	"encoding/json"
	"log"
	"sync"
)

// Структура для хранения id и имен чатов и каналов
var UsersKeywordsChatsMap = make(map[int]map[string]int64)
var muUKC sync.RWMutex // Мьютекс для синхронизации доступа к UsersKeywordsChatsMap

func InitUsersKeywordsChatsMap(userID int) {
	muUKC.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer muUKC.Unlock()
	UsersKeywordsChatsMap[userID] = make(map[string]int64)
}

func AddKeywordToLocalStore(chatID int64, keyword string, userID int) {
	muUKC.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer muUKC.Unlock()
    UsersKeywordsChatsMap[userID][keyword] = chatID
}

func DeleteKeywordFromLocalStore(chatID int64, keyword string, userID int) {
	muUKC.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer muUKC.Unlock()
	if keywordsListForUserId, ok := UsersKeywordsChatsMap[userID]; ok {
		delete(keywordsListForUserId, keyword)
	}
}

func DeleteChatFromLocalStore(chatID int64, userID int) {
	muUKC.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer muUKC.Unlock()
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
	muUKC.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer muUKC.Unlock()    
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

	muUKC.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer muUKC.Unlock()
	for curKeyword, curChatId := range UsersKeywordsChatsMap[userID] {
		if (curChatId == chatID) {
			keywordsList = append(keywordsList, curKeyword)
		}
	}
	return keywordsList
}

func GetUserChatIDForKeyword(keyword string, userID int) (int64) {
	muUKC.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer muUKC.Unlock()
	return UsersKeywordsChatsMap[userID][keyword] 
}

func GetKeywordsMapForUser(userID int) (map[string]int64) {
	muUKC.Lock() // Блокировка мьютекса для исключительного доступа к мапе.
	defer muUKC.Unlock()
	return UsersKeywordsChatsMap[userID]
}