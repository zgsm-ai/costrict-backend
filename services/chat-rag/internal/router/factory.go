package router

import (
	"github.com/zgsm-ai/chat-rag/internal/config"
	"github.com/zgsm-ai/chat-rag/internal/logger"
	ssemantic "github.com/zgsm-ai/chat-rag/internal/router/strategies/semantic"
	"go.uber.org/zap"
)

// NewRunner creates a strategy instance based on config
func NewRunner(cfg config.RouterConfig) Strategy {
	switch cfg.Strategy {
	case "semantic", "":
		return ssemantic.New(cfg.Semantic)
	default:
		logger.Info("semantic router: no strategy matched",
			zap.String("strategy", cfg.Strategy),
		)
		return nil
	}
}
