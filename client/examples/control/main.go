package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"hppr.dev/stoke"
)

func main() {
	ctx := context.Background()
	var err error
	isTest := true

	var keyStore stoke.PublicKeyStore

	if isTest {
		log.Println("SKIPPING AUTH")
		keyStore = stoke.NewTestPublicKeyStore(jwt.New(jwt.SigningMethodNone))
	} else {
		log.Println("USING AUTH")
		keyStore, err = stoke.NewPerRequestPublicKeyStore("http://172.17.0.1:8080/api/pkeys", ctx)
	}

	if err != nil {
		log.Println(err)
		panic("An error occurred while creating public key store")
	}

	log.Println("Ship controller started.")
	
	mux := http.NewServeMux()

	mux.Handle("/location",
		stoke.AuthFunc(
			func (res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(http.StatusOK)
				res.Write([]byte(`{location: "Alpha Quadrant"}`))
			},
			keyStore,
			stoke.RequireToken().WithClaimListPart("ctl", "nav"),
		),
	)
	mux.Handle("/speed",
		stoke.AuthFunc(
			func (res http.ResponseWriter, req *http.Request) {
				// TODO connect this to engine
				res.WriteHeader(http.StatusOK)
				res.Write([]byte(`{speed: "Warp 9"}`))
			},
			keyStore,
			stoke.RequireToken().WithClaimListPart("ctl", "sp"),
		),
	)
	mux.Handle(
		"/my-token",
		stoke.AuthFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write([]byte(fmt.Sprintf("I got token: %s", stoke.Token(rq.Context()).Raw)))
			},
			keyStore,
			stoke.RequireToken(),
		),
	)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Printf("Listening returned an error: %v", err)
	}

}
