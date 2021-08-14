package raftadapter

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/xaraphix/Sif/internal/raft/raftadapter/protos"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type RaftGRPCServer struct {
	pb.UnimplementedRaftRPCAdapterServer
}

type RaftRPCServer struct {
}

func (s *RaftRPCServer) Start() {

	address := "0.0.0.0:50051"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	fmt.Printf("Server is listening on %v ...", address)

	server := grpc.NewServer()
	pb.RegisterRaftRPCAdapterServer(server, &RaftGRPCServer{})

	server.Serve(lis)
}

func (s *RaftGRPCServer) RequestVoteFromPeer(ctx context.Context, vr *pb.VoteRequest) (*pb.VoteResponse, error) {

	response := &pb.VoteResponse{
		PeerId:      int32(3),
		Term:        int32(32),
		VoteGranted: true,
	}

	return response, nil
}

func (s *RaftGRPCServer) ReplicateLog(ctx context.Context, vr *pb.LogRequest) (*pb.LogResponse, error) {

	response := &pb.LogResponse{
		FollowerId: int32(3),
		Term:       int32(32),
		AckLength:  int32(32),
		Success:    true,
	}

	return response, nil
}

func (s *RaftGRPCServer) BroadcastMessage(ctx context.Context, msg *structpb.Struct) (*pb.BroadcastMessageResponse, error) {

	response := &pb.BroadcastMessageResponse{}
	return response, nil
}
