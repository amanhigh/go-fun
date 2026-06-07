package audit

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// PluginRegistry registers and resolves audit plugins.
// Used by the manager to build the catalog and dispatch execution.
type PluginRegistry struct {
	plugins map[barkat.AuditID]Plugin
}

// NewPluginRegistry creates an empty registry.
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{plugins: make(map[barkat.AuditID]Plugin)}
}

// RegisterPlugin adds a plugin. Returns error if ID or order conflicts.
func (r *PluginRegistry) RegisterPlugin(p Plugin) error {
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
func (r *PluginRegistry) GetPlugin(id barkat.AuditID) (Plugin, common.HttpError) {
	p, exists := r.plugins[id]
	if !exists {
		return nil, common.NewHttpError("Audit not found", http.StatusNotFound)
	}
	return p, nil
}

// ListCatalog returns audit metadata sorted by order.
func (r *PluginRegistry) ListCatalog() []barkat.Audit {
	plugins := make([]Plugin, 0, len(r.plugins))
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
