package key

import "time"

type KeyMeta struct {
	Expires time.Time
}

func (m *KeyMeta) SetExpires(t time.Time) {
	m.Expires = t
}

func (m *KeyMeta) ExpiresAt() time.Time {
	return m.Expires
}
