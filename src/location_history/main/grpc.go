package main

import (
	pb "common/protobuff"
	"common/utils"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)



type server struct {
    pb.UnimplementedLocationHistoryServiceServer
}


func startGRPC(){
	lis, err := net.Listen("tcp", ":" + GRPC_PORT)
    if err != nil {
        log.Printf("failed to listen: %v\n", err)
    }
    s := grpc.NewServer()
    pb.RegisterLocationHistoryServiceServer(s, &server{})
	log.Printf("server listening at %v\n", lis.Addr())
    
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err)
    }
}

func (s *server) UpdateHistory(ctx context.Context, req *pb.LocationUpdateRequest) (*pb.LocationUpdateReply, error) {
    username := req.GetUsername()
    longitude := req.GetLongitude()
    latitude := req.GetLatitude()

    if err := utils.CheckUsername(username); err != nil {
        return &pb.LocationUpdateReply{Status: pb.Status_FAILED,Error: err.Error()}, err 
    }

    if err := utils.CheckCoordinates(longitude,latitude); err != nil {
        return &pb.LocationUpdateReply{Status: pb.Status_FAILED, Error: err.Error()}, err 
    }

    if err := updateHistoryByUsername(username,longitude,latitude); err != nil {
        return &pb.LocationUpdateReply{Status: pb.Status_FAILED, Error: err.Error()}, err  
    }
    
    return &pb.LocationUpdateReply{Status: pb.Status_SUCCESS,Error: ""}, nil
}
