package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"hppr.dev/stoke"
)

func main() {
	ctx := context.Background()
	var err error
	isTest := false

	var keyStore stoke.PublicKeyStore

	if isTest {
		fmt.Println("SKIPPING AUTH")
		keyStore, err = stoke.NewTestPublicKeyStore(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMiIsIm5hbWUiOiJKIiwiaWF0IjoxNTE2MjM5MDIyfQ.vQIsenNApq3yE85JJXf72i2zz60r50mNbbFJ-7r_Bv4",
		)
	} else {
		fmt.Println("USING AUTH")
		keyStore, err = stoke.NewPerRequestPublicKeyStore("http://localhost:8080/api/pkeys", ctx)
	}

	if err != nil {
		fmt.Print(err)
		panic("An error occurred while creating public key store")
	}

	log.Println("Starting engines...")
	
	mux := http.NewServeMux()

	mux.Handle("/speed",
		stoke.AuthFunc(
			func (res http.ResponseWriter, req *http.Request) {
				log.Println("Requested speed.")
				res.Write([]byte("To Warp 9!"))
			},
			keyStore,
			stoke.RequireToken().WithClaim("role", "eng"),
		),
	)
	mux.Handle(
		"/my-token",
		stoke.AuthFunc(
			func(rs http.ResponseWriter, rq *http.Request) { rs.Write([]byte(fmt.Sprintf("I got token: %s", stoke.Token(rq.Context()).Raw))) },
			keyStore,
			stoke.RequireToken(),
		),
	)

	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Printf("Listening returned an error: %v", err)
	}

}
