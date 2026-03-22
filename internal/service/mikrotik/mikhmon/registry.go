package mikhmon

import (
	"mikmongo/internal/service"
)

type Registry struct {
	Voucher   *MikhmonVoucherService
	Profile   *MikhmonProfileService
	Report    *MikhmonReportService
	Generator *MikhmonGeneratorService
	Expire    *MikhmonExpireService
}

func NewRegistry(routerSvc *service.RouterService) *Registry {
	generator := NewMikhmonGeneratorService()
	return &Registry{
		Voucher:   NewMikhmonVoucherService(routerSvc, generator),
		Profile:   NewMikhmonProfileService(routerSvc),
		Report:    NewMikhmonReportService(routerSvc),
		Generator: generator,
		Expire:    NewMikhmonExpireService(routerSvc),
	}
}
