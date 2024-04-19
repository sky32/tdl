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
			return watch.GetList(logger.Named(cmd.Context(), "ls"))
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
			return tRun(cmd.Context(), func(ctx context.Context, c *telegram.Client, kvd kv.KV) error {
				return watch.Add(logger.Named(ctx, "add"), c, kvd, opts)
			}, limiter)
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
			return tRun(cmd.Context(), func(ctx context.Context, c *telegram.Client, kvd kv.KV) error {
				return watch.Run(logger.Named(ctx, "run"), c, kvd)
			}, limiter)
		},
	}

	return cmd
}
