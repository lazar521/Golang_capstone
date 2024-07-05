package main

import (
	"errors"

	pb "common/protobuff"
	"context"
	"time"

	"google.golang.org/grpc"
)



func notifyLocationHistoryService(username string, xcoord float64, ycoord float64) error {
	conn, err := grpc.Dial(GRPC_HOST + ":" + GRPC_PORT, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        return err
    }
    defer conn.Close()
    c := pb.NewLocationHistoryServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    
	res, err := c.UpdateHistory(ctx, &pb.LocationUpdateRequest{Username: username, Xcoord: xcoord, Ycoord: ycoord})
    if err != nil {
		return err
    }

	if res.Error != "" {
		return errors.New(res.Error)
	}

	return nil
}