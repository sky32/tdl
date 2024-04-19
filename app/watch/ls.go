package watch

import (
	"context"
	"fmt"
	"github.com/iyear/tdl/pkg/config"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/mattn/go-runewidth"
	"go.uber.org/zap"
	"strings"
)

//go:generate go-enum --names --values --flag --nocase

type AddOptions struct {
	ChatId string
}

func GetList(ctx context.Context) error {
	log := logger.From(ctx)
	printTable(config.Config.Watch.Chats)
	log.Info(
		"Watch Chats:",
		zap.Any("chats", config.Config.Watch.Chats),
	)

	return nil
}

func printTable(result []config.WatchChat) {
	fmt.Printf("\nWatch Chats:\n%s %s %s\n",
		trunc("ID", 15),
		trunc("Name", 30),
		trunc("PreTemplate", 100),
	)
	for _, r := range result {
		fmt.Printf("%s %s %s\n",
			trunc(r.Chat, 15),
			trunc(r.Name, 30),
			trunc(r.PreTemplate, 100),
		)
	}
}

func trunc(s string, len int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "-"
	}

	return runewidth.FillRight(runewidth.Truncate(s, len, "..."), len)
}
