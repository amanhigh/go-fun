package audit

import (
	"fmt"
	"sort"

	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// AuditPluginRegistry registers and resolves audit plugins.
// Used by the manager to build the catalog and dispatch execution.
type AuditPluginRegistry struct {
	plugins map[barkat.AuditID]AuditPlugin
}

// NewAuditPluginRegistry creates an empty registry.
func NewAuditPluginRegistry() *AuditPluginRegistry {
	return &AuditPluginRegistry{plugins: make(map[barkat.AuditID]AuditPlugin)}
}

// RegisterPlugin adds a plugin. Returns error if ID or order conflicts.
func (r *AuditPluginRegistry) RegisterPlugin(p AuditPlugin) error {
	id := p.ID()
	if _, exists := r.plugins[id]; exists {
		return fmt.Errorf("audit plugin %q is already registered", id)
	}
	order := p.Order()
	for _, existing := range r.plugins {
		if existing.Order() == order {
			return fmt.Errorf("duplicate audit plugin order %d: %q conflicts with %q", order, id, existing.ID())
		}
	}
	r.plugins[id] = p
	return nil
}

// GetPlugin returns the plugin for id, or a 404 HttpError if not found.
func (r *AuditPluginRegistry) GetPlugin(id barkat.AuditID) (AuditPlugin, common.HttpError) {
	p, exists := r.plugins[id]
	if !exists {
		return nil, common.ErrNotFound
	}
	return p, nil
}

// ListCatalog returns audit metadata sorted by order.
func (r *AuditPluginRegistry) ListCatalog() []barkat.Audit {
	plugins := make([]AuditPlugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Order() < plugins[j].Order()
	})
	result := make([]barkat.Audit, len(plugins))
	for i, p := range plugins {
		result[i] = Metadata(p)
	}
	return result
}
