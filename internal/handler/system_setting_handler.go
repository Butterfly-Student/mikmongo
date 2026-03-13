package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/pkg/response"
)

// SystemSettingHandler handles system setting HTTP requests
type SystemSettingHandler struct {
	repo repository.SystemSettingRepository
}

// NewSystemSettingHandler creates a new system setting handler
func NewSystemSettingHandler(repo repository.SystemSettingRepository) *SystemSettingHandler {
	return &SystemSettingHandler{repo: repo}
}

// List handles listing system settings
func (h *SystemSettingHandler) List(c *gin.Context) {
	group := c.Query("group")
	var settings []model.SystemSetting
	var err error
	if group != "" {
		settings, err = h.repo.ListByGroup(c.Request.Context(), group)
	} else {
		settings, err = h.repo.List(c.Request.Context(), 1000, 0)
	}
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, settings)
}

// Get handles getting a system setting by ID
func (h *SystemSettingHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	setting, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, setting)
}

// Upsert handles creating or updating a system setting
func (h *SystemSettingHandler) Upsert(c *gin.Context) {
	var setting model.SystemSetting
	if err := c.ShouldBindJSON(&setting); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.repo.Upsert(c.Request.Context(), &setting); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, setting)
}
