package internal

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	ListenAddr string
	ConfigPath string
}

// AppConfig holds the global application configuration
var AppConfig = &Config{}

func InitFlags(rootCmd *cobra.Command) error {
	// Add command line flags
	rootCmd.Flags().StringVarP(&AppConfig.ListenAddr, "listen", "l", "", "Server listen address (e.g. :8080)")
	rootCmd.Flags().StringVarP(&AppConfig.ConfigPath, "config", "c", "", "Configuration file path")

	return nil
}

// LoadConfig loads configuration from file and environment variables
// @returns {error} Error if configuration loading fails
// @description
// - Loads configuration from config.yaml file
// - Merges environment variables
// - Sets default values for missing configurations
// @throws
// - Configuration file not found error
// - Configuration parsing error
func LoadConfig(configPath string) error {
	// If custom config path is provided, use it
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./data")
		viper.AddConfigPath("./config")
	}

	// Set default values
	viper.SetDefault("server.listen", ":8080")
	viper.SetDefault("database.dsn", "./data/client-manager.db")
	viper.SetDefault("log.level", "info")

	// Enable environment variable override
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create default config
			return nil
		}
		return err
	}

	return nil
}

// ApplyConfig applies command line overrides to the configuration
func ApplyConfig() {
	// Override listen address from command line if provided
	if AppConfig.ListenAddr != "" {
		viper.Set("server.listen", AppConfig.ListenAddr)
	}
}

func GetListenAddr() string {
	port := viper.GetString("server.listen")
	if port == "" {
		port = ":8080"
	}
	return port
}
