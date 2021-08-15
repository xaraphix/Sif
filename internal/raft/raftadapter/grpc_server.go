package raftadapter

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/xaraphix/Sif/internal/raft"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCServer struct {
	pb.UnimplementedRaftServer
	raftnode *raft.RaftNode
}

func NewGRPCServer(rn *raft.RaftNode) GRPCServer {
	s := GRPCServer{}
	s.raftnode = rn
	return s
}

func (s *GRPCServer) Start() {

	address := "0.0.0.0:50051"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	fmt.Printf("Server is listening on %v ...", address)
	server := grpc.NewServer()
	pb.RegisterRaftServer(server, &GRPCServer{})
	server.Serve(lis)
}

func (s *GRPCServer) RequestVoteFromPeer(ctx context.Context, vr *pb.VoteRequest) (*pb.VoteResponse, error) {
	return s.raftnode.ElectionMgr.GetResponseForVoteRequest(s.raftnode, vr)
}

func (s *GRPCServer) ReplicateLog(ctx context.Context, vr *pb.LogRequest) (*pb.LogResponse, error) {
	return s.raftnode.LogMgr.RespondToLogReplicationRequest(s.raftnode, vr)
}

func (s *GRPCServer) BroadcastMessage(ctx context.Context, msg *structpb.Struct) (*pb.BroadcastMessageResponse, error) {
	return s.raftnode.LogMgr.RespondToBroadcastMsgRequest(s.raftnode, msg)
}