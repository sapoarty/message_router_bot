package structures

import (
	"sync"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MediaGroups это карта, которая хранит группы медиафайлов для каждого userID.
var MediaGroups = make(map[int]*UserMediaGroups)
var muMGS sync.RWMutex


func InitMediaGroupsForUser(userID int) {
	muMGS.Lock()
	defer muMGS.Unlock()
	MediaGroups[userID] = &UserMediaGroups{
		Groups: make(map[string]*MediaGroup),
	}
}

func InitMediaGroupsForGroupID(userID int, groupID string) {
	muMGS.Lock()
	defer muMGS.Unlock()
	MediaGroups[userID].Groups[groupID] = &MediaGroup{
		Messages: make([]*tgbotapi.Message, 0),
	}
}

func AddMediaGroupsMessages(userID int, groupID string, message *tgbotapi.Message) {
	muMGS.Lock()
	defer muMGS.Unlock()

	if _, userExists := MediaGroups[userID]; !userExists {
		MediaGroups[userID] = &UserMediaGroups{
			Groups: make(map[string]*MediaGroup),
		}
	}

	if mediaGroup, groupExists := MediaGroups[userID].Groups[groupID]; groupExists {
		mediaGroup.Messages = append(mediaGroup.Messages, message)
	} else {
		MediaGroups[userID].Groups[groupID] = &MediaGroup{
			Messages: []*tgbotapi.Message{message},
		}
	}
}

func ClearMediaGroupMessages(userID int, groupID string) {
	muMGS.Lock()
	defer muMGS.Unlock()
	if userGroups, userExists := MediaGroups[userID]; userExists {
		if mediaGroup, groupExists := userGroups.Groups[groupID]; groupExists {
			mediaGroup.Messages = nil
			mediaGroup.Timer.Stop() // Не забудьте остановить таймер перед удалением
			delete(userGroups.Groups, groupID)
		}
		if len(userGroups.Groups) == 0 {
			delete(MediaGroups, userID)
		}
	}
}