package mikrotik

import (
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/pkg/mikrotik/script"
)

// ScriptService provides script generation operations
// Note: This service doesn't require router connection as it generates scripts locally
type ScriptService struct {
	onLoginGen *script.OnLoginGenerator
}

// NewScriptService creates a new Script service
func NewScriptService() *ScriptService {
	return &ScriptService{
		onLoginGen: script.NewOnLoginGenerator(),
	}
}

// GenerateOnLoginScript generates an on-login script for user profile
func (s *ScriptService) GenerateOnLoginScript(req *domain.ProfileRequest) string {
	return s.onLoginGen.Generate(req)
}

// ParseOnLoginScript parses Mikhmon metadata from an existing on-login script
func (s *ScriptService) ParseOnLoginScript(scriptStr string) *domain.ProfileRequest {
	return s.onLoginGen.Parse(scriptStr)
}

// GenerateExpiredAction generates the action script for when user expires
func (s *ScriptService) GenerateExpiredAction(expireMode string) string {
	return s.onLoginGen.GenerateExpiredAction(expireMode)
}

// GenerateExpireMonitorScript generates the global scheduler script
func (s *ScriptService) GenerateExpireMonitorScript() string {
	return s.onLoginGen.GenerateExpireMonitorScript()
}
