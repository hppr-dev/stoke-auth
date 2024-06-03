package main

import (
	"context"
	"engine/proto"
	"log"
	"net"
	"slices"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"hppr.dev/stoke"
)

type engineServer struct {
	proto.UnimplementedEngineRoomServer
	rpm         float32
	speed       float32
	temperature float32
	fuel        float32
	
	commands   chan int
	statusSubs []chan statusReport
	mutex      sync.RWMutex
}

type statusReport struct {
	statusType proto.StatusType
	value float32
}

func (e *engineServer) Start(ctx context.Context) {
	e.rpm = 0
	e.temperature = 80
	e.fuel = 100
	e.speed = 0
	targetSpeed := 0

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Engine stopped")
				return

			case cmd := <-e.commands:
				log.Printf("Engine received target command: %d", cmd)
				targetSpeed += cmd

			case <- time.After(time.Second):
				diff := float32(targetSpeed) - e.speed / 13
				if diff > 0 {
					e.updateRPM(diff)
					e.updateTemp(diff * 0.05)
					e.updateSpeed()
				}
				if e.rpm > 0 {
					e.updateFuel()
				}

			}
		}
	}()
}

func (e *engineServer) updateRPM(inc float32) {
	e.rpm += inc
	e.sendUpdate(statusReport{ statusType : proto.StatusType_RPM, value: e.rpm } )
}

func (e *engineServer) updateTemp(inc float32) {
	e.temperature += inc
	e.sendUpdate(statusReport{ statusType : proto.StatusType_TEMPERATURE, value: e.temperature } )
}

func (e *engineServer) updateSpeed() {
	e.speed = e.rpm / 13
	e.sendUpdate(statusReport{ statusType : proto.StatusType_SPEED, value: e.speed } )
}

func (e *engineServer) updateFuel() {
	e.fuel -= e.rpm * 0.05
	e.sendUpdate(statusReport{ statusType : proto.StatusType_FUEL, value: e.fuel } )
}

func (e *engineServer) sendUpdate(report statusReport) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	if len(e.statusSubs) > 0 {
		for _, c := range e.statusSubs {
			select {
			case c <- report:
				// We don't want to block if the channel is full
			default:
				log.Println("Unable to send update, buffer full")
			}
		}
	}
}

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

func main() {
	log.Print("Starting server on port 6060....")

	listener, err := net.Listen("tcp", "0.0.0.0:6060")
	if err != nil {
		log.Fatalf("Could not listen on 0.0.0.0:6060 -- %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	keyStore, err := stoke.NewPerRequestPublicKeyStore("http://172.17.0.1:8080/api/pkeys", ctx)
	if err != nil {
		log.Fatalf("Could not get stoke public keys: %v", err)
	}
	s := grpc.NewServer(
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
