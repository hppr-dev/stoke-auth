package cfg

import "context"

// Cluster holds HA/cluster options. When Enabled is true, key persistence is disabled
// and /api/pkeys returns a merged JWKS from all discovered peers.
type Cluster struct {
	Enabled     bool     `json:"enabled"`
	Discovery   string   `json:"discovery"`    // "static" (default) or "k8s" (future)
	StaticPeers []string `json:"static_peers"` // base URLs, e.g. https://stoke-1:8080
	RefreshSec  int      `json:"refresh_sec"`  // seconds between peer refresh; default 30
	// InstanceID is a unique identifier for this replica (e.g. "stoke1", "stoke2"). When set,
	// signing key kids are prefixed so merged JWKS from multiple replicas keeps all keys distinct.
	InstanceID string `json:"instance_id"`
}

type clusterCtxKey struct{}

func (c Cluster) withContext(ctx context.Context) context.Context {
	c2 := c
	if c2.RefreshSec <= 0 {
		c2.RefreshSec = 30
	}
	return context.WithValue(ctx, clusterCtxKey{}, &c2)
}

// ClusterFromContext returns the Cluster from ctx, or nil if not set.
func ClusterFromContext(ctx context.Context) *Cluster {
	v := ctx.Value(clusterCtxKey{})
	if v == nil {
		return nil
	}
	if c, ok := v.(*Cluster); ok {
		return c
	}
	return nil
}
