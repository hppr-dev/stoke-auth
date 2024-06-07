package main

import (
	"context"
	"engine/proto"
	"log"
	"sync"
	"time"
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

