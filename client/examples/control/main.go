package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"engine/proto"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hppr.dev/stoke"
)

func main() {
	ctx := context.Background()
	var err error
	isTest := os.Getenv("STOKE_TEST") == "yes"
	engineURL := os.Getenv("ENGINE_URL")

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
				ctx := req.Context()
				// Forward received token to grpc
				client, err := grpc.NewClient(engineURL,
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)
				if err != nil {
					log.Printf("An error occurred creating engine grpc client: %v", err)
					res.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				engineRoom := proto.NewEngineRoomClient(client)
				reply, err := engineRoom.SpeedCommand(ctx,
					&proto.SpeedRequest{
						Direction: proto.SpeedCommandDirection_UP,
						Increment: 10,
					},
					stoke.Credentials().Token(stoke.Token(ctx).Raw).DisableSecurity().CallOption(),
				)
				if err != nil {
					log.Printf("An error occurred calling engine grpc: %v", err)
					res.WriteHeader(http.StatusInternalServerError)
					return
				}

				res.WriteHeader(http.StatusOK)
				res.Write([]byte(fmt.Sprintf(`{message: "%s"}`, reply.Response)))
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
