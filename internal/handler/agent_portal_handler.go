package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"mikmongo/internal/repository"
	"mikmongo/internal/service"
	"mikmongo/pkg/jwt"
	"mikmongo/pkg/response"
)

// AgentPortalHandler handles agent self-service portal endpoints
type AgentPortalHandler struct {
	agentRepo  repository.SalesAgentRepository
	invoiceSvc *service.AgentInvoiceService
	saleSvc    *service.HotspotSaleService
	jwtService *jwt.Service
}

// NewAgentPortalHandler creates a new agent portal handler
func NewAgentPortalHandler(
	agentRepo repository.SalesAgentRepository,
	invoiceSvc *service.AgentInvoiceService,
	saleSvc *service.HotspotSaleService,
	jwtService *jwt.Service,
) *AgentPortalHandler {
	return &AgentPortalHandler{
		agentRepo:  agentRepo,
		invoiceSvc: invoiceSvc,
		saleSvc:    saleSvc,
		jwtService: jwtService,
	}
}

// Login authenticates an agent for portal access
func (h *AgentPortalHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	agent, err := h.agentRepo.GetByUsername(c.Request.Context(), req.Username)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(agent.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if agent.Status != "active" {
		response.Forbidden(c, "agent account is inactive")
		return
	}

	token, err := h.jwtService.GenerateAgent(agent.ID, agent.Username)
	if err != nil {
		response.InternalServerError(c, "failed to generate token")
		return
	}

	agent.PasswordHash = ""
	response.OK(c, gin.H{"token": token, "agent": agent})
}

// GetProfile returns the agent's profile
func (h *AgentPortalHandler) GetProfile(c *gin.Context) {
	agentID := c.MustGet("agent_id").(string)
	id, _ := uuid.Parse(agentID)
	agent, err := h.agentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "agent not found")
		return
	}
	agent.PasswordHash = ""
	response.OK(c, agent)
}

// ChangePassword changes the agent's password
func (h *AgentPortalHandler) ChangePassword(c *gin.Context) {
	agentID := c.MustGet("agent_id").(string)
	id, _ := uuid.Parse(agentID)

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	agent, err := h.agentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "agent not found")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalServerError(c, "failed to hash password")
		return
	}
	agent.PasswordHash = string(hash)

	if err := h.agentRepo.Update(c.Request.Context(), agent); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "password changed"})
}

// GetInvoices returns the agent's invoices
func (h *AgentPortalHandler) GetInvoices(c *gin.Context) {
	agentID := c.MustGet("agent_id").(string)
	id, _ := uuid.Parse(agentID)

	limit, offset := getPagination(c)
	invs, count, err := h.invoiceSvc.ListByAgent(c.Request.Context(), id, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, invs, &response.Meta{Total: count, Limit: limit, Offset: offset})
}

// GetInvoice returns a specific invoice (verified ownership)
func (h *AgentPortalHandler) GetInvoice(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	inv, err := h.invoiceSvc.GetInvoice(c.Request.Context(), invoiceID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	agentID := c.MustGet("agent_id").(string)
	if inv.AgentID != agentID {
		response.Forbidden(c, "access denied")
		return
	}
	response.OK(c, inv)
}

// RequestPayment submits payment proof for an invoice, transitioning it to "review" status
func (h *AgentPortalHandler) RequestPayment(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	agentID := c.MustGet("agent_id").(string)

	var req struct {
		PaidAmount float64 `json:"paid_amount" binding:"required,min=0"`
		Notes      string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	inv, err := h.invoiceSvc.RequestPayment(c.Request.Context(), invoiceID, agentID, req.PaidAmount, req.Notes)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, inv)
}

// GetSales returns the agent's voucher sales history
func (h *AgentPortalHandler) GetSales(c *gin.Context) {
	agentID := c.MustGet("agent_id").(string)
	id, _ := uuid.Parse(agentID)

	filter := repository.HotspotSaleFilter{
		SalesAgentID: &id,
	}

	limit, offset := getPagination(c)
	sales, count, err := h.saleSvc.ListSales(c.Request.Context(), filter, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, sales, &response.Meta{Total: count, Limit: limit, Offset: offset})
}
