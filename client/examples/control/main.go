package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"engine/proto"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/net/websocket"
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

	grpcClient, err := grpc.NewClient(engineURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println(err)
		panic("An error occured while connecting to engine")
	}
	defer grpcClient.Close()

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

				engineRoom := proto.NewEngineRoomClient(grpcClient)

				// Forward received token to grpc
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
	mux.HandleFunc("/foobar", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{ 
			Handler: websocket.Handler(grpcFoobarWebsocket(keyStore, grpcClient)),
		}
		s.ServeHTTP(w, req)
	})

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

func grpcFoobarWebsocket(keyStore stoke.PublicKeyStore, grpcConn *grpc.ClientConn) websocket.Handler {
	return func (ws *websocket.Conn) {
		defer ws.Close()
		req := ws.Request()
		ctx := req.Context()

		rawToken := req.URL.Query().Get("token")
		_, err := keyStore.ParseClaims(ctx, rawToken, stoke.RequireToken().WithClaimListPart("ctl", "acc"))
		if err != nil {
			log.Printf("ws token invalid or not found: %s", err)
			return
		}

		engineRoom := proto.NewEngineRoomClient(grpcConn)

		// Forward token to engine
		fbClient, err := engineRoom.FooBarTest(ctx, stoke.Credentials().Token(rawToken).DisableSecurity().CallOption())
		if err != nil {
			log.Printf("Could not call FooBarTest: %v", err)
			return
		}

		// This would be better suited as Unary, but we want to test streaming.
		for {
			var in []byte
			if err := websocket.Message.Receive(ws, &in); errors.Is(err, io.EOF) {
				log.Printf("Connection closed.")
				return
			} else if err != nil {
				log.Printf("Could not read message: %v", err)
				return
			}

			if err := fbClient.Send(&proto.SimpleMessage{ Message: string(in) }); err != nil {
				log.Printf("error sending grpc: %v", err)
				return
			}

			resp, err := fbClient.Recv()
			if err != nil {
				log.Printf("error sending grpc: %v", err)
				return
			}

			log.Printf("Got response from grpc: %s", resp.Message)

			if err := websocket.Message.Send(ws, resp.Message); err != nil {
				log.Printf("Could not write message: %v", err)
				return
			}
		}
	}
}
