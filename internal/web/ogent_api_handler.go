package web

import (
	"log"
	"net/http"
	"stoke/internal/ctx"
	"stoke/internal/ent/ogent"
)

type EntityAPI struct {
	*ogent.OgentHandler
	Context *ctx.Context
}

func NewEntityAPIHandler(prefix string, context *ctx.Context) http.Handler {
	if len(prefix) > 0 && prefix[len(prefix) - 1] == '/' {
		prefix = prefix[0:len(prefix) - 1]
	}
	hdlr, err := ogent.NewServer(
		EntityAPI{
			OgentHandler: ogent.NewOgentHandler(context.DB),
			Context: context,
		},
		ogent.WithPathPrefix(prefix),
	)
	if err != nil {
		log.Printf("An error occured creating an entity handler: %v", err)
		return nil
	}
	return hdlr

}
