package usr

type Provider interface {
	Init() error
	ValidateUser(user, pass string) bool
}

type MultiProvider struct {
	providers []Provider
}

func (m MultiProvider) Add(p Provider) {
	m.providers = append(m.providers, p)
}

func (m MultiProvider) Init() error {
	for _, p := range m.providers {
		err := p.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m MultiProvider) ValidateUser(user, pass string) bool {
	for _, p := range m.providers {
		if p.ValidateUser(user, pass) {
			return true
		}
	}
	return false
}
