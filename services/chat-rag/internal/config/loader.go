package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/zgsm-ai/chat-rag/internal/logger"
	"go.uber.org/zap"
)

// LoadYAML loads yaml from the specified file path using viper
func LoadYAML[T any](path string) (*T, error) {
	var yaml T

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read YAML: %w", err)
	}

	if err := viper.Unmarshal(&yaml); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &yaml, nil
}

// MustLoadConfig loads configuration and panics if there's an error
func MustLoadConfig(configPath string) Config {
	c, err := LoadYAML[Config](configPath)
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Apply defaults: if fallbackModelName not set, use the first candidate
	if c != nil && c.Router.Semantic.Routing.FallbackModelName == "" {
		if len(c.Router.Semantic.Routing.Candidates) > 0 {
			c.Router.Semantic.Routing.FallbackModelName = c.Router.Semantic.Routing.Candidates[0].ModelName
		}
	}

	// Align rule engine prefix defaults with original plugin logic
	if c != nil {
		if c.Router.Semantic.RuleEngine.BodyPrefix == "" {
			c.Router.Semantic.RuleEngine.BodyPrefix = "body."
		}
		if c.Router.Semantic.RuleEngine.HeaderPrefix == "" {
			c.Router.Semantic.RuleEngine.HeaderPrefix = "header."
		}
	}

	// Align stripCodeFences default behavior with plugin:
	// default to true when the key is not explicitly set in YAML
	if c != nil {
		// inputExtraction.protocol default
		if c.Router.Semantic.InputExtraction.Protocol == "" {
			c.Router.Semantic.InputExtraction.Protocol = "openai"
		}
		// inputExtraction.userJoinSep default
		if c.Router.Semantic.InputExtraction.UserJoinSep == "" {
			c.Router.Semantic.InputExtraction.UserJoinSep = "\n\n"
		}
		// inputExtraction.stripCodeFences default (only when key not set)
		if !viper.IsSet("router.semantic.inputExtraction.stripCodeFences") {
			c.Router.Semantic.InputExtraction.StripCodeFences = true
		}
		// inputExtraction.codeFenceRegex default is empty string (no-op) — keep if unset
		// inputExtraction.maxUserMessages default
		if !viper.IsSet("router.semantic.inputExtraction.maxUserMessages") || c.Router.Semantic.InputExtraction.MaxUserMessages == 0 {
			c.Router.Semantic.InputExtraction.MaxUserMessages = 100
		}
		// inputExtraction.maxHistoryBytes default
		if !viper.IsSet("router.semantic.inputExtraction.maxHistoryBytes") || c.Router.Semantic.InputExtraction.MaxHistoryBytes == 0 {
			c.Router.Semantic.InputExtraction.MaxHistoryBytes = 2048
		}
		// inputExtraction.maxHistoryMessages default
		if !viper.IsSet("router.semantic.inputExtraction.maxHistoryMessages") || c.Router.Semantic.InputExtraction.MaxHistoryMessages == 0 {
			c.Router.Semantic.InputExtraction.MaxHistoryMessages = 5
		}
		// Apply idle timeout defaults
		if c.LLMTimeout.IdleTimeoutMs <= 0 {
			c.LLMTimeout.IdleTimeoutMs = 30000
			logger.Info("llm idle timeout not set, using default", zap.Int("idleTimeoutMs", c.LLMTimeout.IdleTimeoutMs))
		}
		if c.LLMTimeout.TotalIdleTimeoutMs <= 0 {
			c.LLMTimeout.TotalIdleTimeoutMs = 30000
			logger.Info("llm total idle timeout not set, using default", zap.Int("totalIdleTimeoutMs", c.LLMTimeout.TotalIdleTimeoutMs))
		}
	}

	// Apply forward configuration defaults
	if c != nil {
		// forward.enabled default
		if !viper.IsSet("forward.enabled") {
			c.Forward.Enabled = false
		}
		// forward.defaultTarget default
		if !viper.IsSet("forward.defaultTarget") {
			c.Forward.DefaultTarget = ""
		}
	}

	logger.Info("loaded config", zap.Any("config", c))
	return *c
}

// LoadRulesConfig loads the rules configuration from etc/rules.yaml
func LoadRulesConfig() (*RulesConfig, error) {
	// Get the project root directory path
	projectRoot, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Build the configuration file path
	configPath := filepath.Join(projectRoot, "etc", "rules.yaml")

	// Load YAML configuration using config package
	rulesConfig, err := LoadYAML[RulesConfig](configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load rules config: %w", err)
	}

	// Log successful loading
	logger.Info("Rules configuration loaded successfully at service startup",
		zap.String("config_path", configPath),
		zap.Int("agents_count", len(rulesConfig.Agents)))

	return rulesConfig, nil
}
