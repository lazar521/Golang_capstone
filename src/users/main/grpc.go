package main

import (
	"context"
	"errors"
	"time"

	pb "common/protobuff" // Importing the protobuf generated code

	"google.golang.org/grpc"
)

// notifyLocationHistoryService sends location updates to the gRPC location history service
var notifyLocationHistoryService = func(username string, longitude float64, latitude float64) error {
	// Dial the gRPC server
	conn, err := grpc.Dial(GRPC_HOST+":"+GRPC_PORT, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create a new client for the LocationHistoryService
	c := pb.NewLocationHistoryServiceClient(conn)

	// Set a timeout for the context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Send the UpdateHistory request with the username, longitude, and latitude
	res, err := c.UpdateHistory(ctx, &pb.LocationUpdateRequest{Username: username, Longitude: longitude, Latitude: latitude})
	if err != nil {
		return err
	}

	// Check if the response contains an error message
	if res.Error != "" {
		return errors.New(res.Error)
	}

	// Return nil if no error occurred
	return nil
}
