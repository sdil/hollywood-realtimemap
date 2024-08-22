package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func createVehicleHandler(engine *actor.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vid := r.URL.Query().Get("id")

		pid := engine.Registry.GetPID("video", vid)
		if pid == nil {
			http.Error(w, "Error lookup PID", http.StatusInternalServerError)
			return
		}

		resp := engine.Request(pid, &positionRequest{}, time.Minute)
		result, err := resp.Result()
		if err != nil {
			fmt.Println("Error requesting position", err)
			http.Error(w, "Error requesting position", http.StatusInternalServerError)
			return
		}
		if result, ok := result.(positionResponse); ok {
			fmt.Fprintf(w, "Position: %v", result.Position)
		} else {
			http.Error(w, "Error getting position", http.StatusInternalServerError)
		}
	}
}

func main() {
	ctx := context.Background()

	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		fmt.Println("Error creating actor system", err)
		return
	}

	http.HandleFunc("/vehicle", createVehicleHandler(engine))

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

			pid := engine.Registry.GetPID("video", *vid)
			if pid == nil {
				pid = engine.Spawn(NewVehicle, "video", actor.WithID(*vid))
			}

			engine.Send(pid, &Position{
				Latitude:  *event.VehiclePosition.Latitude,
				Longitude: *event.VehiclePosition.Longitude,
			})
		}
	}, ctx)

	<-ingressDone
}
