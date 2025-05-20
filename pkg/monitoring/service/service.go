package service

import (
	"github.com/nessi-dev/nessi/pkg/monitoring/freshness"
)

// Provider is a service provider for the monitoring system
type Provider struct {
	freshnessManager *freshness.Manager
}

// NewProvider creates a new service provider
func NewProvider() *Provider {
	return &Provider{
		freshnessManager: freshness.NewManager(),
	}
}

// GetFreshnessManager returns the freshness manager
func (p *Provider) GetFreshnessManager() (*freshness.Manager, error) {
	return p.freshnessManager, nil
}
