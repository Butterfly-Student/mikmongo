package mikhmon

import (
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	gorosmikhmon "github.com/Butterfly-Student/go-ros/repository/mikhmon"
)

type MikhmonGeneratorService struct {
	repo gorosmikhmon.GeneratorRepository
}

func NewMikhmonGeneratorService() *MikhmonGeneratorService {
	return &MikhmonGeneratorService{
		repo: gorosmikhmon.NewGeneratorRepository(),
	}
}

func (s *MikhmonGeneratorService) GenerateUsername(config *mikhmonDomain.GeneratorConfig) string {
	return s.repo.GenerateUsername(config)
}

func (s *MikhmonGeneratorService) GeneratePassword(config *mikhmonDomain.GeneratorConfig) string {
	return s.repo.GeneratePassword(config)
}

func (s *MikhmonGeneratorService) GeneratePair(mode string, config *mikhmonDomain.GeneratorConfig) *mikhmonDomain.GeneratorResult {
	return s.repo.GeneratePair(mode, config)
}
