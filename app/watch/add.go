package watch

import (
	"context"
	"github.com/fatih/color"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/peers"
	"github.com/iyear/tdl/pkg/config"
	"github.com/iyear/tdl/pkg/kv"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/iyear/tdl/pkg/storage"
	"github.com/iyear/tdl/pkg/utils"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func Add(ctx context.Context, c *telegram.Client, kvd kv.KV, opts AddOptions) error {
	color.Green("Adding chat\n")
	log := logger.From(ctx)

	var peer peers.Peer
	var err error

	// Log the start of the Add operation
	log.Info("Starting to add chat", zap.String("chat_id", opts.ChatId))

	manager := peers.Options{Storage: storage.NewPeers(kvd)}.Build(c.API())

	// Handle 'me' as a special case
	if opts.ChatId == "me" {
		peer, err = manager.Self(ctx)
		if err != nil {
			log.Error("Failed to get self", zap.Error(err))
			return err
		}
		log.Info("Got self", zap.Any("peer", peer))
	} else {
		peer, err = utils.Telegram.GetInputPeer(ctx, manager, opts.ChatId)
		if err != nil {
			log.Error("Failed to get input peer", zap.String("chat_id", opts.ChatId), zap.Error(err))
			return err
		}
		log.Info("Got input peer", zap.Any("peer", peer), zap.String("chat_id", opts.ChatId))
	}
	// Check if the chat ID already exists in the configuration
	if isChatIdPresent(strconv.FormatInt(peer.ID(), 10), config.Config) {
		log.Warn("Chat ID already exists in the configuration", zap.String("chat_id", opts.ChatId))
		color.Red("chat_id %s already exists\n", opts.ChatId)
		return nil
	}

	// Create a new watch chat configuration
	chat := config.WatchChat{
		Type:          1,
		Chat:          strconv.FormatInt(peer.ID(), 10),
		Thread:        0,
		Input:         []int{},
		Filter:        "true",
		OnlyMedia:     false,
		WithContent:   false,
		Raw:           false,
		All:           false,
		LastMessageAt: "",
		Name:          peer.VisibleName(),
		PreTemplate:   cleanFolderName(peer.VisibleName()),
	}

	// Append the new chat to the configuration and save it
	config.Config.Watch.Chats = append(config.Config.Watch.Chats, chat)
	err = config.SaveConfig()
	if err != nil {
		log.Error("Failed to save configuration", zap.Error(err))
		return err
	}
	log.Info("Added chat to configuration", zap.Any("peer", peer), zap.Any("chat", chat))

	// Log the successful completion of the Add operation
	color.Green("Add chat %s[%s] success...\n", chat.Name, chat.Chat)
	log.Info("Successfully added chat", zap.String("name", chat.Name), zap.String("chat_id", chat.Chat))

	return nil
}

func cleanFolderName(name string) string {
	// 用 _ 替换非法字符
	replacementFunc := func(r rune) rune {
		switch r {
		case ' ', '\\', '/', '<', '>', ':', '"', '|', '?', '*':
			return '_'
		default:
			return r
		}
	}
	return strings.Map(replacementFunc, name)
}

func isChatIdPresent(chatId string, config *config.Configuration) bool {
	for _, chat := range config.Watch.Chats {
		if chat.Chat == chatId {
			return true
		}
	}
	return false
}
