package main

import (
	pb "common/protobuff"
	"common/utils"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

// server struct implements the gRPC service interface defined in the protobuf
type server struct {
	pb.UnimplementedLocationHistoryServiceServer
}

// startGRPC initializes and starts the gRPC server
func startGRPC() {
	// Listen on the specified TCP port
	lis, err := net.Listen("tcp", ":"+GRPC_PORT)
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the LocationHistoryServiceServer with the gRPC server
	pb.RegisterLocationHistoryServiceServer(s, &server{})
	log.Printf("server listening at %v\n", lis.Addr())

	// Serve gRPC server
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}

// UpdateHistory handles the UpdateHistory RPC call
// It updates the location history for a given username and coordinates
func (s *server) UpdateHistory(ctx context.Context, req *pb.LocationUpdateRequest) (*pb.LocationUpdateReply, error) {
	username := req.GetUsername()
	longitude := req.GetLongitude()
	latitude := req.GetLatitude()

	// Validate the username
	if err := utils.CheckUsername(username); err != nil {
		return &pb.LocationUpdateReply{Status: pb.Status_FAILED, Error: err.Error()}, err
	}

	// Validate the coordinates
	if err := utils.CheckCoordinates(longitude, latitude); err != nil {
		return &pb.LocationUpdateReply{Status: pb.Status_FAILED, Error: err.Error()}, err
	}

	// Update the location history in the database
	if err := updateHistoryByUsername(username, longitude, latitude); err != nil {
		return &pb.LocationUpdateReply{Status: pb.Status_FAILED, Error: err.Error()}, err
	}

	// Return a successful response
	return &pb.LocationUpdateReply{Status: pb.Status_SUCCESS, Error: ""}, nil
}
