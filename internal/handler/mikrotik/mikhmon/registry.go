package mikhmon

import (
	"mikmongo/internal/service/mikrotik"
)

type Registry struct {
	Voucher *VoucherHandler
	Profile *ProfileHandler
	Report  *ReportHandler
	Expire  *ExpireHandler
}

func NewRegistry(mkRegistry *mikrotik.Registry) *Registry {
	return &Registry{
		Voucher: NewVoucherHandler(mkRegistry.Mikhmon.Voucher, mkRegistry.Mikhmon.Generator),
		Profile: NewProfileHandler(mkRegistry.Mikhmon.Profile),
		Report:  NewReportHandler(mkRegistry.Mikhmon.Report),
		Expire:  NewExpireHandler(mkRegistry.Mikhmon.Expire),
	}
}
