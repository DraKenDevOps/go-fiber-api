package handlers

import (
	"go-fiber-api/config"

	"gorm.io/gorm"
)

type ApiHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewApiHandler(cfg *config.Config, db *gorm.DB) *ApiHandler {
	return &ApiHandler{cfg: cfg, db: db}
}
