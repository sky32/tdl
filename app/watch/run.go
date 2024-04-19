package watch

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/gotd/td/telegram"
	"github.com/iyear/tdl/app/chat"
	"github.com/iyear/tdl/app/dl"
	config "github.com/iyear/tdl/pkg/config"
	"github.com/iyear/tdl/pkg/kv"
	"github.com/iyear/tdl/pkg/logger"
	"go.uber.org/zap"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func Run(ctx context.Context, c *telegram.Client, kvd kv.KV) error {
	log := logger.From(ctx)
	exportConf := config.Config.Watch.Export
	if exportConf.Dir == "" {
		exportConf.Dir = "downloads"
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGQUIT)
	go func() {
		<-signalCh
		color.Red("Received SIGQUIT, exiting...")
		cancel()
	}()
	color.Green("Start watching...\n")
	log.Info("Start watching...")
	for {
		for i, watchChat := range config.Config.Watch.Chats {
			log.Info("Watching chat", zap.String("chat", watchChat.Name))
			originInput := watchChat.Input
			if err := processChatInput(&watchChat); err != nil {
				return err
			}
			output := filepath.Join(exportConf.Dir, watchChat.PreTemplate, "list.json")
			dir := filepath.Dir(output)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
			options, err := export(ctx, c, kvd, watchChat, output, log)
			config.Config.Watch.Chats[i].LastMessageAt = time.Now().Format("2006-01-02 15:04:05")
			config.Config.Watch.Chats[i].Input = originInput
			if err := config.SaveConfig(); err != nil {
				return err
			}
			if err := download(ctx, c, kvd, log, exportConf, err, watchChat, options); err != nil {
				return err
			}
		}
		duration := time.Duration(config.Config.Watch.Interval) * time.Minute
		color.Yellow("\nWaiting for %s to continue...\n", duration)
		select {
		case <-time.After(duration):
		case <-ctx.Done():
			color.Yellow("Exiting...")
			return nil
		}
	}
	//goland:noinspection GoUnreachableCode
	return nil
}

func download(ctx context.Context, c *telegram.Client, kvd kv.KV, log *zap.Logger, exportConf config.WatchExport, err error, watchChat config.WatchChat, options chat.ExportOptions) error {
	color.Green("\nDownload chat\n")
	opts := dl.Options{
		Dir:        exportConf.Dir,
		RewriteExt: exportConf.RewriteExt,
		SkipSame:   exportConf.SkipSame,
		Template:   filepath.Join(watchChat.PreTemplate, exportConf.Template),
		Files:      []string{options.Output},
		Include:    exportConf.Include,
		Exclude:    exportConf.Exclude,
		Desc:       exportConf.Desc,
		Takeout:    exportConf.Takeout,
		Restart:    exportConf.Restart,
		Continue:   exportConf.Continue,
	}
	log.Info("Downloading files", zap.Any("opts", opts))
	err = dl.Run(ctx, c, kvd, opts)
	log.Info("Downloaded files", zap.Any("opts", opts), zap.Any("err", err))
	if err != nil {
		return err
	}
	return nil
}

func export(ctx context.Context, c *telegram.Client, kvd kv.KV, watchChat config.WatchChat, output string, log *zap.Logger) (chat.ExportOptions, error) {
	color.Green("\nExporting chat\n")
	opts := chat.ExportOptions{
		Type:        watchChat.Type,
		Chat:        watchChat.Chat,
		Thread:      watchChat.Thread,
		Input:       watchChat.Input,
		Output:      output,
		Filter:      watchChat.Filter,
		OnlyMedia:   watchChat.OnlyMedia,
		WithContent: watchChat.WithContent,
		Raw:         watchChat.Raw,
		All:         watchChat.All,
	}
	log.Info("Exporting chat", zap.Any("opts", opts))
	err := chat.Export(ctx, c, kvd, opts)
	log.Info("Exported chat", zap.Any("opts", opts), zap.Any("err", err))
	return opts, err
}

func processChatInput(watchChat *config.WatchChat) error {
	if watchChat.LastId > 0 && len(watchChat.Input) == 0 {
		watchChat.Input = []int{watchChat.LastId + 1}
	}
	switch watchChat.Type {
	case chat.ExportTypeTime, chat.ExportTypeId:
		// set default value
		switch len(watchChat.Input) {
		case 0:
			watchChat.Input = []int{0, math.MaxInt}
		case 1:
			watchChat.Input = append(watchChat.Input, math.MaxInt)
		}

		if len(watchChat.Input) != 2 {
			return fmt.Errorf("input data should be 2 integers when watchChat type is %s", watchChat.Type)
		}

		// sort helper
		if watchChat.Input[0] > watchChat.Input[1] {
			watchChat.Input[0], watchChat.Input[1] = watchChat.Input[1], watchChat.Input[0]
		}
	case chat.ExportTypeLast:
		if len(watchChat.Input) != 1 {
			return fmt.Errorf("input data should be 1 integer when watchChat type is %s", watchChat.Type)
		}
	default:
		return fmt.Errorf("unknown watchChat type: %s", watchChat.Type)
	}
	return nil
}
