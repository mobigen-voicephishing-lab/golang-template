package handler

import (
	"github.com/labstack/echo/v5"

	"github.com/mobigen/golang-web-template/internal/adapter/inbound/http/dto"
	"github.com/mobigen/golang-web-template/internal/infrastructure/config"
)

// VersionHandler version endpoint handler
type VersionHandler struct{}

// NewVersionHandler create VersionHandler instance
func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

// GetVersion return app version
// @Summary Get Server Version
// @Description get server version info
// @Tags version
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.HTTPResponse[dto.VersionResponse] "app info(name, version, hash)"
// @Router /version [get]
func (h *VersionHandler) GetVersion(c *echo.Context) error {
	return OK(c, dto.VersionResponse{
		Name:      config.Name,
		Version:   config.Version,
		BuildHash: config.BuildHash,
	})
}
