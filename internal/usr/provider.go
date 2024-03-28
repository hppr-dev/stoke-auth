package usr

import "stoke/internal/ent"

type Provider interface {
	Init() error
	GetUserClaims(user, pass string) (ent.Claims, error)
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

func (m MultiProvider) GetUserClaims(username, password string) (ent.Claims, error) {
	var claims ent.Claims
	for _, p := range m.providers {
		provClaims, _ := p.GetUserClaims(username, password)
		claims = append(claims, provClaims...)
	}
	if len(claims) == 0 {
		return nil, NoClaimsFound{}
	}
	return claims, nil
}

type NoClaimsFound struct {
}

func (NoClaimsFound) Error() string {
	return "No claims found for user"
}
