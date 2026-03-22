package mikhmon

import (
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

type GenerateVoucherRequest struct {
	Quantity   int    `json:"quantity" binding:"required,min=1,max=1000"`
	Server     string `json:"server,omitempty"`
	Profile    string `json:"profile" binding:"required"`
	Mode       string `json:"mode" binding:"required,oneof=vc up"`
	NameLength int    `json:"name_length" binding:"min=3,max=12"`
	Prefix     string `json:"prefix,omitempty"`
	CharSet    string `json:"char_set" binding:"required"`
	TimeLimit  string `json:"time_limit,omitempty"`
	DataLimit  string `json:"data_limit,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

type VoucherBatchResponse struct {
	Code      string            `json:"code"`
	Quantity  int               `json:"quantity"`
	Profile   string            `json:"profile"`
	Server    string            `json:"server"`
	TimeLimit string            `json:"time_limit,omitempty"`
	DataLimit string            `json:"data_limit,omitempty"`
	Vouchers  []VoucherResponse `json:"vouchers"`
}

type VoucherResponse struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	Profile  string `json:"profile,omitempty"`
	Server   string `json:"server,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Code     string `json:"code,omitempty"`
	Mode     string `json:"mode,omitempty"`
	Date     string `json:"date,omitempty"`
}

type GetVouchersQuery struct {
	Server  string `form:"server"`
	Profile string `form:"profile"`
	Mode    string `form:"mode"`
	Limit   int    `form:"limit"`
}

func convertVoucherBatchToResponse(batch *mikhmonDomain.VoucherBatch) VoucherBatchResponse {
	vouchers := make([]VoucherResponse, len(batch.Vouchers))
	for i, v := range batch.Vouchers {
		vouchers[i] = convertVoucherToResponse(&v)
	}
	return VoucherBatchResponse{
		Code:      batch.Code,
		Quantity:  batch.Quantity,
		Profile:   batch.Profile,
		Server:    batch.Server,
		TimeLimit: batch.TimeLimit,
		DataLimit: batch.DataLimit,
		Vouchers:  vouchers,
	}
}

func convertVoucherToResponse(v *mikhmonDomain.Voucher) VoucherResponse {
	return VoucherResponse{
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
