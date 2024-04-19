package dl

import (
	"github.com/iyear/tdl/pkg/config"
	"github.com/iyear/tdl/pkg/downloader"
	"strconv"
)

func WatchBegin(elem downloader.Elem) {
	e := elem.(*iterElem)
	chatId := e.from.ID()
	messageId := e.fromMsg.ID
	config.Config.Watch.Mu.RLock()
	var found bool
	for i, chat := range config.Config.Watch.Chats {
		if chat.Chat == strconv.FormatInt(chatId, 10) {
			found = true
			config.Config.Watch.Mu.RUnlock()
			config.Config.Watch.Mu.Lock()
			config.Config.Watch.Chats[i].HandleIds = append(config.Config.Watch.Chats[i].HandleIds, messageId)
			break
		}
	}
	if found {
		defer config.Config.Watch.Mu.Unlock()
	} else {
		defer config.Config.Watch.Mu.Unlock()
	}
}

func WatchEnd(elem downloader.Elem) {
	e := elem.(*iterElem)
	chatId := e.from.ID()
	messageId := e.fromMsg.ID
	var found bool
	config.Config.Watch.Mu.RLock()
	for i, chat := range config.Config.Watch.Chats {
		if chat.Chat == strconv.FormatInt(chatId, 10) {
			found = true
			config.Config.Watch.Mu.RUnlock()
			config.Config.Watch.Mu.Lock()
			for j, handledId := range chat.HandleIds {
				if handledId == messageId {
					config.Config.Watch.Chats[i].HandleIds = append(
						chat.HandleIds[:j],
						chat.HandleIds[j+1:]...,
					)
					break
				}
			}
			minId := messageId
			if len(chat.HandleIds) > 0 {
				minId = chat.HandleIds[0]
				for _, id := range chat.HandleIds {
					if id < minId {
						minId = id
					}
				}
			}
			if chat.LastId < messageId && minId == messageId {
				config.Config.Watch.Chats[i].LastId = messageId
				_ = config.SaveConfig()
			}
			break
		}
	}
	if found {
		defer config.Config.Watch.Mu.Unlock()
	} else {
		defer config.Config.Watch.Mu.Unlock()
	}
}
