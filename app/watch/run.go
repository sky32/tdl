package watch

import (
	"context"
	"fmt"
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
	for {
		log := logger.From(ctx)
		log.Info("Starting periodic run")
		exportConf := config.Config.Watch.Export
		if exportConf.Dir == "" {
			exportConf.Dir = "downloads"
		}
		for _, watchChat := range config.Config.Watch.Chats {
			originInput := watchChat.Input
			if watchChat.LastId > 0 && len(watchChat.Input) == 0 {
				watchChat.Input = []int{watchChat.LastId}
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

			output := filepath.Join(exportConf.Dir, watchChat.PreTemplate, "list.json")
			fmt.Printf("Exporting chat %s to %s\n", watchChat.Chat, output)
			dir := filepath.Dir(output)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}

			options := chat.ExportOptions{
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
			log.Info("Exporting chat", zap.Any("options", options))
			err := chat.Export(ctx, c, kvd, options)
			if err != nil {
				return err
			}
			watchChat.LastMessageAt = time.Now().Format("2006-01-02 15:04:05")
			watchChat.Input = originInput
			err = config.SaveConfig()
			if err != nil {
				return err
			}
			log.Info("Downloading files", zap.Any("options", exportConf))
			err = dl.Run(ctx, c, kvd, dl.Options{
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
			})
			if err != nil {
				return err
			}
			log.Info("Finished periodic run")
		}
		os.Exit(0)
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGQUIT)

		go func() {
			for {
				duration := time.Duration(config.Config.Watch.Interval) * time.Minute
				fmt.Printf("Sleeping for %d minutes\n", duration/time.Minute)
				select {
				case <-time.After(duration):
					fmt.Println("Finished sleeping")
				case <-sigCh:
					fmt.Println("Interrupted during sleep")
					os.Exit(0)
				}
			}
		}()

		fmt.Println("Waiting for interrupt...")

		<-sigCh
	}
}
