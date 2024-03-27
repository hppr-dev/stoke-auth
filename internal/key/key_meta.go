package key

import "time"

type KeyMeta struct {
	Expires time.Time
	Renews  time.Time
}

func (m *KeyMeta) SetExpires(t time.Time) {
	m.Expires = t
}

func (m *KeyMeta) ExpiresAt() time.Time {
	return m.Expires
}

func (m *KeyMeta) SetRenews(t time.Time) {
	m.Renews = t
}

func (m *KeyMeta) RenewsAt() time.Time {
	return m.Renews
}
