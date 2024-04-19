package config

import (
	"fmt"
	"github.com/iyear/tdl/app/chat"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

var Config *Configuration

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

type Configuration struct {
	Watch struct {
		Chats    []WatchChat `yaml:"chats"`
		Export   WatchExport `yaml:"export"`
		Interval int         `yaml:"interval"`
		Mu       sync.RWMutex
	} `yaml:"watch"`
}

func LoadConfig() error {
	configFileName := "config.yaml"
	configDir := "."
	configFilePath := filepath.Join(configDir, configFileName)

	if _, err := os.Stat(configFilePath); err != nil {
		if os.IsNotExist(err) {
			_, _ = os.Create(configFilePath)
			fmt.Printf("Config file %s not found, creating a new one...\n", configFilePath)
		}
	}

	viper.SetConfigName(filepath.Base(configFileName))
	viper.AddConfigPath(configDir)
	viper.SetConfigType(filepath.Ext(configFileName)[1:])

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

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	var decodeHook mapstructure.DecodeHookFunc = func(from, to reflect.Type, v interface{}) (interface{}, error) {
		if from.Kind() == reflect.String && to == reflect.TypeOf(time.Time{}) {
			return time.Parse(time.RFC3339, v.(string))
		}
		return v, nil
	}

	if err := viper.Unmarshal(&Config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		decodeHook,
	))); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	viper.WatchConfig()

	return nil
}

func SaveConfig() error {
	viper.Set("watch", Config.Watch)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
