package main

import (
	"context"
	"fmt"
	"os"
	"github.com/anthdm/hollywood/actor"
	"net/http"
)

var pids = make(map[string]*actor.PID)

func main() {
	ctx := context.Background()

	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		fmt.Println("Error creating actor system", err)
		return
	}

	// http.HandleFunc("/vehicle", createVehicleHandler(actorSystem))

	fmt.Println("Server is starting on port 8080...")
	go func() {
		host := "localhost:8080"
		if os.Getenv("RENDER") == "true" {
			host = "0.0.0.0:10000"
		}
		err = http.ListenAndServe(host, nil)
		if err != nil {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()

	ingressDone := ConsumeVehicleEvents(func(event *Event) {
		if event.VehiclePosition.HasValidPosition() {
			vid := &event.VehicleId

			pid := pids[*vid]
			if pid == nil {
				pid = engine.Spawn(NewVehicle, *vid)
				pids[*vid] = pid
			}

			engine.Send(pid, &Position{
				Latitude:  *event.VehiclePosition.Latitude,
				Longitude: *event.VehiclePosition.Longitude,
			})
		}
	}, ctx)

	<-ingressDone
}
