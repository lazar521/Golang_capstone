package main

import (
	"errors"

	pb "common/protobuff"
	"context"
	"time"

	"google.golang.org/grpc"
)



func notifyLocationHistoryService(username string, longitude float64, latitude float64) error {
	conn, err := grpc.Dial(GRPC_HOST + ":" + GRPC_PORT, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        return err
    }
    defer conn.Close()
    c := pb.NewLocationHistoryServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    
	res, err := c.UpdateHistory(ctx, &pb.LocationUpdateRequest{Username: username, Longitude: longitude, Latitude: latitude})
    if err != nil {
		return err
    }

	if res.Error != "" {
		return errors.New(res.Error)
	}

	return nil
}