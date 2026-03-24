package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// CashManagementHandler handles cash entry and petty cash HTTP endpoints.
type CashManagementHandler struct {
	svc *service.CashManagementService
}

// NewCashManagementHandler creates a new CashManagementHandler.
func NewCashManagementHandler(svc *service.CashManagementService) *CashManagementHandler {
	return &CashManagementHandler{svc: svc}
}

// --- Cash Entry endpoints ---

type createCashEntryRequest struct {
	Type            string  `json:"type" binding:"required,oneof=income expense"`
	Source          string  `json:"source" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	Description     string  `json:"description" binding:"required"`
	PaymentMethod   string  `json:"payment_method" binding:"required"`
	BankName        *string `json:"bank_name"`
	AccountNumber   *string `json:"account_number"`
	PettyCashFundID *string `json:"petty_cash_fund_id"`
	EntryDate       string  `json:"entry_date"` // YYYY-MM-DD, defaults to now
	Notes           *string `json:"notes"`
}

// ListEntries handles GET /cash-entries with filters.
func (h *CashManagementHandler) ListEntries(c *gin.Context) {
	filter := repository.CashEntryFilter{
		Type:   c.Query("type"),
		Source: c.Query("source"),
		Status: c.Query("status"),
	}
	if v := c.Query("date_from"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			response.BadRequest(c, "invalid date_from, use YYYY-MM-DD")
			return
		}
		filter.DateFrom = &t
	}
	if v := c.Query("date_to"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			response.BadRequest(c, "invalid date_to, use YYYY-MM-DD")
			return
		}
		filter.DateTo = &t
	}
	if v := c.Query("created_by"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			response.BadRequest(c, "invalid created_by")
			return
		}
		filter.CreatedBy = &id
	}
	if v := c.Query("petty_cash_fund_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			response.BadRequest(c, "invalid petty_cash_fund_id")
			return
		}
		filter.PettyCashFundID = &id
	}

	limit, offset := getPagination(c)
	entries, count, err := h.svc.ListEntries(c.Request.Context(), filter, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, entries, &response.Meta{Total: count, Limit: limit, Offset: offset})
}

// CreateEntry handles POST /cash-entries.
func (h *CashManagementHandler) CreateEntry(c *gin.Context) {
	var req createCashEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	createdBy, ok := userID.(string)
	if !ok || createdBy == "" {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	entryDate := time.Now()
	if req.EntryDate != "" {
		t, err := time.Parse("2006-01-02", req.EntryDate)
		if err != nil {
			response.BadRequest(c, "invalid entry_date, use YYYY-MM-DD")
			return
		}
		entryDate = t
	}

	entry := &model.CashEntry{
		Type:            req.Type,
		Source:          req.Source,
		Amount:          req.Amount,
		Description:     req.Description,
		PaymentMethod:   req.PaymentMethod,
		BankName:        req.BankName,
		AccountNumber:   req.AccountNumber,
		PettyCashFundID: req.PettyCashFundID,
		EntryDate:       entryDate,
		CreatedBy:       createdBy,
		Notes:           req.Notes,
	}

	if err := h.svc.CreateEntry(c.Request.Context(), entry); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, entry)
}

// GetEntry handles GET /cash-entries/:id.
func (h *CashManagementHandler) GetEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	entry, err := h.svc.GetEntry(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, entry)
}

// UpdateEntry handles PUT /cash-entries/:id.
func (h *CashManagementHandler) UpdateEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Prevent updating protected fields
	for _, key := range []string{"id", "entry_number", "status", "created_by", "approved_by", "approved_at", "created_at"} {
		delete(updates, key)
	}

	entry, err := h.svc.UpdateEntry(c.Request.Context(), id, updates)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, entry)
}

// DeleteEntry handles DELETE /cash-entries/:id.
func (h *CashManagementHandler) DeleteEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.DeleteEntry(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}

// ApproveEntry handles POST /cash-entries/:id/approve.
func (h *CashManagementHandler) ApproveEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	userID, _ := c.Get("user_id")
	approvedBy, ok := userID.(string)
	if !ok || approvedBy == "" {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	entry, err := h.svc.Approve(c.Request.Context(), id, approvedBy)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, entry)
}

// RejectEntry handles POST /cash-entries/:id/reject.
func (h *CashManagementHandler) RejectEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	entry, err := h.svc.Reject(c.Request.Context(), id, req.Reason)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, entry)
}

// --- Petty Cash Fund endpoints ---

type createFundRequest struct {
	FundName       string  `json:"fund_name" binding:"required"`
	InitialBalance float64 `json:"initial_balance" binding:"required,gte=0"`
	CustodianID    string  `json:"custodian_id" binding:"required,uuid"`
}

type updateFundRequest struct {
	FundName *string `json:"fund_name"`
	Status   *string `json:"status"`
}

type topUpFundRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// ListFunds handles GET /petty-cash.
func (h *CashManagementHandler) ListFunds(c *gin.Context) {
	limit, offset := getPagination(c)
	funds, count, err := h.svc.ListFunds(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, funds, &response.Meta{Total: count, Limit: limit, Offset: offset})
}

// CreateFund handles POST /petty-cash.
func (h *CashManagementHandler) CreateFund(c *gin.Context) {
	var req createFundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	fund := &model.PettyCashFund{
		FundName:       req.FundName,
		InitialBalance: req.InitialBalance,
		CustodianID:    req.CustodianID,
		Status:         "active",
	}
	if err := h.svc.CreateFund(c.Request.Context(), fund); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, fund)
}

// GetFund handles GET /petty-cash/:id.
func (h *CashManagementHandler) GetFund(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	fund, err := h.svc.GetFund(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, fund)
}

// UpdateFund handles PUT /petty-cash/:id.
func (h *CashManagementHandler) UpdateFund(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req updateFundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	fund, err := h.svc.GetFund(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	if req.FundName != nil {
		fund.FundName = *req.FundName
	}
	if req.Status != nil {
		fund.Status = *req.Status
	}
	fund.UpdatedAt = time.Now()

	if err := h.svc.UpdateFund(c.Request.Context(), fund); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, fund)
}

// TopUpFund handles POST /petty-cash/:id/topup.
func (h *CashManagementHandler) TopUpFund(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req topUpFundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.TopUpFund(c.Request.Context(), id, req.Amount); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	fund, err := h.svc.GetFund(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, fund)
}

// --- Report endpoints ---

// GetCashFlow handles GET /reports/cash-flow?from=YYYY-MM-DD&to=YYYY-MM-DD.
func (h *CashManagementHandler) GetCashFlow(c *gin.Context) {
	from, err := time.Parse("2006-01-02", c.Query("from"))
	if err != nil {
		response.BadRequest(c, "invalid 'from' date, use YYYY-MM-DD")
		return
	}
	to, err := time.Parse("2006-01-02", c.Query("to"))
	if err != nil {
		response.BadRequest(c, "invalid 'to' date, use YYYY-MM-DD")
		return
	}

	report, err := h.svc.GetCashFlow(c.Request.Context(), from, to)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, report)
}

// GetCashBalance handles GET /reports/cash-balance.
func (h *CashManagementHandler) GetCashBalance(c *gin.Context) {
	report, err := h.svc.GetCashBalance(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, report)
}

// GetReconciliation handles GET /reports/reconciliation?from=YYYY-MM-DD&to=YYYY-MM-DD.
func (h *CashManagementHandler) GetReconciliation(c *gin.Context) {
	from, err := time.Parse("2006-01-02", c.Query("from"))
	if err != nil {
		response.BadRequest(c, "invalid 'from' date, use YYYY-MM-DD")
		return
	}
	to, err := time.Parse("2006-01-02", c.Query("to"))
	if err != nil {
		response.BadRequest(c, "invalid 'to' date, use YYYY-MM-DD")
		return
	}

	report, err := h.svc.GetReconciliation(c.Request.Context(), from, to)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, report)
}
