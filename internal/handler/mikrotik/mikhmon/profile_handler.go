package mikhmon

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"mikmongo/internal/dto/mikrotik/mikhmon"
	mikhmonSvc "mikmongo/internal/service/mikrotik/mikhmon"
	"mikmongo/pkg/response"
)

type ProfileHandler struct {
	profileSvc *mikhmonSvc.MikhmonProfileService
}

func NewProfileHandler(profileSvc *mikhmonSvc.MikhmonProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileSvc: profileSvc,
	}
}

func (h *ProfileHandler) Create(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	var req mikhmon.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.profileSvc.CreateProfile(c.Request.Context(), routerID, req.ToDomain()); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, gin.H{"message": "profile created"})
}

func (h *ProfileHandler) Update(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	profileID := c.Param("id")
	if profileID == "" {
		response.BadRequest(c, "invalid profile id")
		return
	}

	var req mikhmon.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.profileSvc.UpdateProfile(c.Request.Context(), routerID, profileID, req.ToDomain(profileID)); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "profile updated"})
}

func (h *ProfileHandler) GenerateScript(c *gin.Context) {
	var req mikhmon.GenerateScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	script := h.profileSvc.GenerateOnLoginScript(req.ToDomain())
	response.OK(c, mikhmon.OnLoginScriptResponse{Script: script})
}
