package main

import (
	"context"
	"github.com/anthdm/hollywood/actor"
	"time"
	"fmt"
)

type Vehicle struct {
	id       string
	position []Position
}

type Position struct {
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

func NewVehicle() actor.Receiver {
	return &Vehicle{}
}

func (v *Vehicle) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		v.id = ctx.PID().ID
		fmt.Println("actor started", v.id)
	case actor.Stopped:
		fmt.Println("actor stopped")
	case *Position:
		v.position = append(v.position, *msg)
		fmt.Println(v.position)
		// fmt.Println("actor has received", v.id, msg.Latitude, msg.Longitude)
	}
}

func (v *Vehicle) PostStop(ctx context.Context) error {
	// Do nothing
	return nil
}
