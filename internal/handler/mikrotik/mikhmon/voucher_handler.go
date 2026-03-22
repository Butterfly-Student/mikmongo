package mikhmon

import (
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mikhdto "mikmongo/internal/dto/mikrotik/mikhmon"
	mikhmonservice "mikmongo/internal/service/mikrotik/mikhmon"
	"mikmongo/pkg/response"
)

type VoucherHandler struct {
	voucherSvc   *mikhmonservice.MikhmonVoucherService
	generatorSvc *mikhmonservice.MikhmonGeneratorService
}

func NewVoucherHandler(voucherSvc *mikhmonservice.MikhmonVoucherService, generatorSvc *mikhmonservice.MikhmonGeneratorService) *VoucherHandler {
	return &VoucherHandler{
		voucherSvc:   voucherSvc,
		generatorSvc: generatorSvc,
	}
}

func (h *VoucherHandler) GenerateBatch(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	var req mikhdto.GenerateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	voucherReq := &mikhmonDomain.VoucherGenerateRequest{
		Quantity:   req.Quantity,
		Server:     req.Server,
		Profile:    req.Profile,
		Mode:       req.Mode,
		NameLength: req.NameLength,
		Prefix:     req.Prefix,
		CharSet:    req.CharSet,
		TimeLimit:  req.TimeLimit,
		DataLimit:  req.DataLimit,
		Comment:    req.Comment,
	}

	batch, err := h.voucherSvc.GenerateBatch(c.Request.Context(), routerID, voucherReq)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, convertVoucherBatchToResponse(batch))
}

func (h *VoucherHandler) GetVouchers(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	comment := c.Query("comment")
	code := c.Query("code")

	var vouchers []*mikhmonDomain.Voucher
	if comment != "" {
		vouchers, err = h.voucherSvc.GetVouchersByComment(c.Request.Context(), routerID, comment)
	} else if code != "" {
		vouchers, err = h.voucherSvc.GetVouchersByCode(c.Request.Context(), routerID, code)
	} else {
		response.BadRequest(c, "either comment or code query parameter is required")
		return
	}

	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, convertVouchersToResponse(vouchers))
}

func (h *VoucherHandler) RemoveBatch(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	comment := c.Query("comment")
	if comment == "" {
		response.BadRequest(c, "comment query parameter is required")
		return
	}

	if err := h.voucherSvc.RemoveVoucherBatch(c.Request.Context(), routerID, comment); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "vouchers removed"})
}

func convertVoucherBatchToResponse(batch *mikhmonDomain.VoucherBatch) mikhdto.VoucherBatchResponse {
	vouchers := make([]mikhdto.VoucherResponse, len(batch.Vouchers))
	for i, v := range batch.Vouchers {
		vouchers[i] = convertVoucherToResponse(&v)
	}
	return mikhdto.VoucherBatchResponse{
		Code:      batch.Code,
		Quantity:  batch.Quantity,
		Profile:   batch.Profile,
		Server:    batch.Server,
		TimeLimit: batch.TimeLimit,
		DataLimit: batch.DataLimit,
		Vouchers:  vouchers,
	}
}

func convertVoucherToResponse(v *mikhmonDomain.Voucher) mikhdto.VoucherResponse {
	return mikhdto.VoucherResponse{
		ID:       v.ID,
		Name:     v.Name,
		Password: v.Password,
		Profile:  v.Profile,
		Server:   v.Server,
		Comment:  v.Comment,
		Code:     v.Code,
		Mode:     v.Mode,
		Date:     v.Date,
	}
}

func convertVouchersToResponse(vouchers []*mikhmonDomain.Voucher) []mikhdto.VoucherResponse {
	result := make([]mikhdto.VoucherResponse, len(vouchers))
	for i, v := range vouchers {
		result[i] = convertVoucherToResponse(v)
	}
	return result
}
