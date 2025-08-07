package middleware

import (
	"ka-cache/config"
	"ka-cache/logger"
)

type MiddlewareManager struct {
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

func NewMiddlewareManager(cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, origins: origins, logger: logger}
}
