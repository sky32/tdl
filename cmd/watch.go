package cmd

import (
	"context"
	"github.com/gotd/td/telegram"
	"github.com/iyear/tdl/app/watch"
	"github.com/iyear/tdl/pkg/kv"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/spf13/cobra"
)

func NewWatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "A set of watch tools",
	}

	cmd.AddCommand(GetWatchList(), AddWatch(), RunWatch())

	return cmd
}

func GetWatchList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "Get your watched chats",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.Named(cmd.Context(), "ls")
			err := watch.GetList(log)
			if err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}

func AddWatch() *cobra.Command {
	var opts watch.AddOptions
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add new watch to chats",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.Named(cmd.Context(), "add")
			err := tRun(cmd.Context(), func(ctx context.Context, c *telegram.Client, kvd kv.KV) error {
				return watch.Add(log, c, kvd, opts)
			}, limiter)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.ChatId, "ChatId", "i", "", "watch chat by id")

	return cmd
}

func RunWatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run watch download",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.Named(cmd.Context(), "run")
			err := tRun(cmd.Context(), func(ctx context.Context, c *telegram.Client, kvd kv.KV) error {
				return watch.Run(log, c, kvd)
			}, limiter)
			if err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}
