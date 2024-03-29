package web

import (
	"log"
	"net/http"
	"stoke/internal/ctx"
)

type PkeyApiHandler struct {
	Context *ctx.Context
}

func (p PkeyApiHandler) ServeHTTP(res http.ResponseWriter, _ *http.Request) {
	b, err := p.Context.Issuer.PublicKeys()
	if err != nil {
		log.Printf("Could not get public keys: %v", err)
	}
	res.Write(b)
}
