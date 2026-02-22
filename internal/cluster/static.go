package cluster

import "context"

// StaticDiscoverer returns a fixed list of peer URLs from config.
type StaticDiscoverer struct {
	URLs []string
}

func (s *StaticDiscoverer) Peers(ctx context.Context) ([]string, error) {
	if s == nil {
		return nil, nil
	}
	out := make([]string, len(s.URLs))
	copy(out, s.URLs)
	return out, nil
}
