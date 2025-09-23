package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Gemini code, Quality unknown

const (
	AppName        string = "dribble"
	AppNameUpper   string = "DRIBBLE"
	configFileName string = "config"
)

type Config struct {
	ShowDrivers []string `mapstructure:"show_drivers"`

	Connections struct {
		Servers map[string]struct {
			ServerName string
			Settings   Connection
		} `mapstructure:"servers"`
		Databases map[string]struct {
			DatabaseName string
			Settings     Connection
		} `mapstructure:"databases"`

		DriverDefaults []Connection `mapstructure:"-"`
	} `mapstructure:"connections"`
	Ui struct {
		ShowDetails bool `mapstructure:"show_details"`

		Theme struct {
			Borders struct {
				Enabled bool
				Color   string
			} `mapstructure:"borders"`
			Colors struct {
				Primary    string
				Secondary  string
				Accent     string
				Text       string
				Background string
			} `mapstructure:"colors"`
		} `mapstructure:"theme"`
	} `mapstructure:"ui"`
}

// returns the base directory for application configuration files
// following the XDG Base Directory Specification.
// It checks XDG_CONFIG_HOME first, then falls back to ~/.config.
func getConfigDir() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome != "" {
		return filepath.Join(configHome, AppName), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".config", AppName), nil
}

// LoadConfig initializes Viper and loads configuration.
func loadConfig() (*Config, *viper.Viper, error) {
	v := viper.New()

	// 1. Set config file name (without extension)
	v.SetConfigName("config")
	v.SetConfigType("yaml") // Set the config file type

	// 2. Add config search paths (in order of precedence if multiple files match)
	// Add current directory first for development/testing
	v.AddConfigPath(".")

	// Add XDG config path
	appConfigDir, err := getConfigDir()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get config directory: %w", err)
	}
	v.AddConfigPath(appConfigDir) // e.g., ~/.config/dribble/

	// Add other paths if needed, e.g., system-wide config:
	// v.AddConfigPath("/etc/" + appName)

	// 3. Set default values

	// Shown Drivers
	v.SetDefault("show_drivers", map[string]struct{}{})

	// Connections
	v.SetDefault("connections.servers", map[string]struct{}{})
	v.SetDefault("connections.databases", map[string]struct{}{})
	v.SetDefault("Connections.DriverDefaults", map[string]struct{}{})
	// v.SetDefault("Connections.Saved", map[string]struct{}{})

	// Ui
	v.SetDefault("ui.show_details", false)

	// Ui, Theme, empty means terminal colors
	v.SetDefault("ui.theme.borders.enabled", true)
	v.SetDefault("ui.theme.borders.color", "")

	v.SetDefault("ui.theme.colors.primary", "")
	v.SetDefault("ui.theme.colors.secondary", "")
	v.SetDefault("ui.theme.colors.accent", "")
	v.SetDefault("ui.theme.colors.text", "")
	v.SetDefault("ui.theme.colors.background", "")

	// 4. Enable automatic environment variable binding
	// This will look for environment variables with the prefix
	// e.g., DRIBBLE_DATABASE_HOST will map to database.host
	v.SetEnvPrefix(AppNameUpper) // DRIBBLE
	v.AutomaticEnv()             // Read ENV variables

	// You can also explicitly bind specific environment variables if needed
	// e.BindEnv("database.password", "DB_PASSWORD") // Binds database.password to DB_PASSWORD env var

	// 5. Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if we're happy with defaults/env vars
			logger.Infof("Config file '%s.yaml' not found in search paths. Using defaults/environment variables.\n", configFileName)
			fmt.Printf("Config file '%s.yaml' not found in search paths. Using defaults/environment variables.\n", configFileName)
		} else {
			// Some other error occurred
			logger.ErrorF("Error reading config file: %s", err)
			panic(fmt.Errorf("failed to read config file: %w", err))
		}
	} else {
		// fmt.Printf("Config file loaded from: %s\n", v.ConfigFileUsed())
	}

	// 6. Unmarshal the configuration into your struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, v, nil
}
