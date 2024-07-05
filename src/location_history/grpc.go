package main

import (
	pb "common/protobuff"
	"common/utils"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)



type server struct {
    pb.UnimplementedLocationHistoryServiceServer
}


func startGRPC(){
	lis, err := net.Listen("tcp", ":" + GRPC_PORT)
    if err != nil {
        fmt.Printf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterLocationHistoryServiceServer(s, &server{})
	fmt.Printf("server listening at %v", lis.Addr())
    
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
    }
}

func (s *server) UpdateHistory(ctx context.Context, req *pb.LocationUpdateRequest) (*pb.LocationUpdateReply, error) {
    username := req.GetUsername()
    xcoord := req.GetXcoord()
    ycoord := req.GetYcoord()

    if err := utils.CheckUsername(username); err != nil {
        return &pb.LocationUpdateReply{Status: pb.Status_FAILED,Error: err.Error()}, err 
    }

    if err := utils.CheckCoordinates(xcoord,ycoord); err != nil {
        return &pb.LocationUpdateReply{Status: pb.Status_FAILED, Error: err.Error()}, err 
    }

    if err := updateHistoryByUsername(username,xcoord,ycoord); err != nil {
        return &pb.LocationUpdateReply{Status: pb.Status_FAILED, Error: err.Error()}, err  
    }
    
    return &pb.LocationUpdateReply{Status: pb.Status_SUCCESS,Error: ""}, nil
}
