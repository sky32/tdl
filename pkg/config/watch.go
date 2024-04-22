package config

import (
	"github.com/iyear/tdl/app/chat"
	"github.com/spf13/viper"
)

type WatchChat struct {
	Type        chat.ExportType `yaml:"type"`
	Chat        string          `yaml:"chat"`
	Thread      int             `yaml:"thread"`
	Input       []int           `yaml:"input"`
	Filter      string          `yaml:"filter"`
	OnlyMedia   bool            `yaml:"onlyMedia"`
	WithContent bool            `yaml:"withContent"`
	Raw         bool            `yaml:"raw"`
	All         bool            `yaml:"all"`

	LastMessageAt string `yaml:"lastMessageAt"`
	Name          string `yaml:"name"`
	PreTemplate   string `yaml:"preTemplate"`
	LastId        int    `yaml:"lastId"`
	HandleIds     []int  `yaml:"handleIds"`
}

type WatchExport struct {
	Dir        string   `yaml:"dir"`
	RewriteExt bool     `yaml:"rewriteExt"`
	SkipSame   bool     `yaml:"skipSame"`
	Template   string   `yaml:"template"`
	URLs       []string `yaml:"URLs"`
	Files      []string `yaml:"files"`
	Include    []string `yaml:"include"`
	Exclude    []string `yaml:"exclude"`
	Desc       bool     `yaml:"desc"`
	Takeout    bool     `yaml:"takeout"`
	Continue   bool     `yaml:"continue"`
	Restart    bool     `yaml:"restart"`
}

func setWatchDefault() {
	viper.SetDefault("watch.export.dir", "downloads")
	viper.SetDefault("watch.export.rewriteExt", false)
	viper.SetDefault("watch.export.skipSame", false)
	viper.SetDefault("watch.export.template", "{{ .DialogID }}/{{ formatDate .MessageDate \"2006-01-02\"}}/{{ .MessageID }}_{{ .FileName  }}")
	viper.SetDefault("watch.export.URLs", []string{})
	viper.SetDefault("watch.export.files", []string{})
	viper.SetDefault("watch.export.include", []string{})
	viper.SetDefault("watch.export.exclude", []string{})
	viper.SetDefault("watch.export.desc", false)
	viper.SetDefault("watch.export.takeout", false)
	viper.SetDefault("watch.export.continue", true)
	viper.SetDefault("watch.export.restart", false)
	viper.SetDefault("watch.interval", 10)
}

func watchToSave() {
	viper.Set("watch", Config.Watch)
}
