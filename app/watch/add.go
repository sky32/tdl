package watch

import (
	"context"
	"fmt"
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
	log := logger.From(ctx)

	var peer peers.Peer
	var err error
	manager := peers.Options{Storage: storage.NewPeers(kvd)}.Build(c.API())
	if opts.ChatId == "me" { // defaults to me(saved messages)
		peer, err = manager.Self(ctx)
		log.Info("Get me", zap.Any("peer", peer))
	} else {
		if isChatIdPresent(opts.ChatId, config.Config) {
			return fmt.Errorf("chat_id %s already exists", opts.ChatId)
		}
		peer, err = utils.Telegram.GetInputPeer(ctx, manager, opts.ChatId)
		log.Info("Get chat", zap.Any("peer", peer), zap.Any("chat_id", opts.ChatId))
	}
	if err != nil {
		return err
	}
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
	config.Config.Watch.Chats = append(
		config.Config.Watch.Chats,
		chat,
	)
	err = config.SaveConfig()
	if err != nil {
		return err
	}
	log.Info("Add chat", zap.Any("peer", peer), zap.Any("chat", chat))
	color.Green("Add chat %s[%s] success...\n", chat.Name, chat.Chat)

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
