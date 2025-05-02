package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

// RPCServer represents a server that handles RPC (Remote Procedure Call) requests and processes payload data.
type RPCServer struct{}

// RPCPayload represents the structure of data payload handled in RPC communication.
// Name specifies the name associated with the payload.
// Data contains the actual payload information.
type RPCPayload struct {
	Name string
	Data string
}

// LogInfo handles the insertion of a log entry into the database using the provided RPC payload and returns a response.
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{Name: payload.Name, Data: payload.Data, CreatedAt: time.Now()})

	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	*resp = "Processed payload vis RPC" + payload.Name

	return nil
}
