package main

import (
	"context"
	"engine/proto"
	"errors"
	"io"
	"log"
	"net"
	"slices"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"hppr.dev/stoke"
)

// SpeedCommand implements proto.EngineRoomServer.
func (e *engineServer) SpeedCommand(ctx context.Context, req *proto.SpeedRequest) (*proto.SpeedReply, error) {
	cmd := int(req.Increment)
	if req.Direction == proto.SpeedCommandDirection_DOWN {
			cmd *= -1
	}
	e.commands <- cmd
	return &proto.SpeedReply{
		Response: "Speed Sent",
	}, nil
}

// StatusStream implements proto.EngineRoomServer.
func (e *engineServer) StatusStream(settings *proto.StatusSettings, server proto.EngineRoom_StatusStreamServer) error {
	statusSub := make(chan statusReport)
	e.mutex.Lock()
	subNum := len(e.statusSubs)
	e.statusSubs = append(e.statusSubs, statusSub)
	e.mutex.Unlock()

	defer func() {
		e.mutex.Lock()
		e.statusSubs = slices.Delete(e.statusSubs, subNum, subNum)
		e.mutex.Unlock()
	}()

	for {
		select {
		case <-server.Context().Done():
			log.Println("Client disconnected.")
			return nil

		case reply := <-statusSub:
			if settings.All ||
					( settings.Rpm && reply.statusType == proto.StatusType_RPM ) ||
					( settings.Fuel && reply.statusType == proto.StatusType_FUEL ) ||
					( settings.Temperature && reply.statusType == proto.StatusType_TEMPERATURE ) ||
					( settings.Speed && reply.statusType == proto.StatusType_SPEED ) {
				server.Send(&proto.StatusReply{
					StatusType: reply.statusType,
					Level:      reply.value,
				})
			}
		}
	}
}

func (e *engineServer) FooBarTest(s proto.EngineRoom_FooBarTestServer) error {
	for {
		msg, err := s.Recv()
		if errors.Is(io.EOF, err) {
			return nil
		}
		if err != nil {
			return err
		}
		switch msg.Message {
		case "foo":
			s.Send( &proto.SimpleMessage{ Message : "bar" } )
		case "hello":
			s.Send( &proto.SimpleMessage{ Message : "world" } )
		default:
			s.Send(msg)
		}
	}
}

func main() {
	log.Print("Starting server on port 6060....")

	listener, err := net.Listen("tcp", "0.0.0.0:6060")
	if err != nil {
		log.Fatalf("Could not listen on 0.0.0.0:6060 -- %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	keyStore, err := stoke.NewPerRequestPublicKeyStore("https://172.17.0.1:8080/api/pkeys", ctx, stoke.ConfigureTLS("/etc/ca.crt"))
	if err != nil {
		log.Fatalf("Could not create stoke public key store: %v", err)
	}

	grpcCreds, err := credentials.NewServerTLSFromFile("/etc/engine.crt", "/etc/engine.key")
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	s := grpc.NewServer(
		grpc.Creds(grpcCreds),
		stoke.NewUnaryTokenInterceptor(keyStore, stoke.RequireToken()).Opt(),
		stoke.NewStreamTokenInterceptor(keyStore, stoke.RequireToken()).Opt(),
	)

	engine := &engineServer{
		commands : make(chan int, 10),
	}

	proto.RegisterEngineRoomServer(s, engine)
	reflection.Register(s)

	engine.Start(ctx)

	s.Serve(listener)

	log.Print("Server Terminated.")
}
