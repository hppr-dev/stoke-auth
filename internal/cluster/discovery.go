package cluster

import "context"

// Discoverer returns the list of peer base URLs (e.g. https://stoke-1:8080) to fetch JWKS from.
type Discoverer interface {
	Peers(ctx context.Context) ([]string, error)
}
