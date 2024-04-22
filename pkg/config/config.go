package config

import (
	"fmt"
	"github.com/iyear/tdl/pkg/consts"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

// Config is a global pointer to the configuration variable, storing the application's configuration information.
var Config *Configuration

// Configuration defines the specific structure of the configuration, including the watchable components.
type Configuration struct {
	Watch struct {
		Chats    []WatchChat  `yaml:"chats"`    // List of chats to monitor
		Export   WatchExport  `yaml:"export"`   // Export settings
		Interval int          `yaml:"interval"` // Interval for monitoring
		Mu       sync.RWMutex // Read-write mutex for concurrent control
	} `yaml:"watch"`
}

// LoadConfig loads the configuration file. It first checks if the configuration file exists; if not, it creates a new one.
// Returns any errors that may occur.
func LoadConfig() error {
	// Initialize Config before using it
	Config = &Configuration{}
	// Step 1: Create or check the configuration file
	configFileName := "config.yaml"
	configDir := consts.DataDir
	configFilePath := filepath.Join(configDir, configFileName)
	if _, err := os.Stat(configFilePath); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(configFilePath); err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}
			// Log the creation instead of printing to stdout
			logCreation(configFilePath)
		}
	}

	// Step 2: Configure Viper and set default values
	return setupViper(configFileName, configDir)
}

// logCreation logs a message indicating the config file has been created.
func logCreation(configFilePath string) {
	// Assuming a simple log function that prints to stderr with timestamps
	fmt.Fprintf(os.Stderr, "[INFO] Config file %s not found, creating a new one...\n", configFilePath)
}

// setupViper configures the Viper instance, setting the configuration file name, path, and type, and sets default values.
func setupViper(configFileName string, configDir string) error {
	// Configure viper and set default values
	viper.SetConfigName(filepath.Base(configFileName))
	viper.AddConfigPath(configDir)
	viper.SetConfigType(filepath.Ext(configFileName)[1:])

	setWatchDefault()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal the configuration using custom decode hooks
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

	return nil
}

// SaveConfig saves the current configuration to a file.
// Returns any errors that may occur.
func SaveConfig() error {
	// Assuming a function watchToSave that updates the Config before saving
	watchToSave()

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
