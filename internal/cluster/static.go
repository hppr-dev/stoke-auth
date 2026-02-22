package cluster

import "context"

// StaticDiscoverer returns a fixed list of peer URLs.
type StaticDiscoverer struct {
	URLs []string
}

// Peers returns a copy of the configured peer URLs.
func (s *StaticDiscoverer) Peers(ctx context.Context) ([]string, error) {
	if s == nil {
		return nil, nil
	}
	out := make([]string, len(s.URLs))
	copy(out, s.URLs)
	return out, nil
}
