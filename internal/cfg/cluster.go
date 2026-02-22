package cfg

import "context"

// Cluster holds HA/cluster options. When Enabled is true, key persistence is disabled
// and /api/pkeys returns a merged JWKS from all discovered peers.
type Cluster struct {
	Enabled     bool     `json:"enabled"`
	Discovery   string   `json:"discovery"`    // "static" (default) or "k8s" (future)
	StaticPeers []string `json:"static_peers"` // base URLs, e.g. https://stoke-1:8080
	RefreshSec  int      `json:"refresh_sec"`  // seconds between peer refresh; default 30
}

type clusterCtxKey struct{}

func (c Cluster) withContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, clusterCtxKey{}, &c)
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
