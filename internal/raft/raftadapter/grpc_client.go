package raftadapter

import (
	"context"
	"log"
	"time"

	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type RaftGRPCClient struct {
	client              pb.RaftClient
	ctx                 context.Context
}

func NewRaftGRPCClient(address string, timeoutIn time.Duration) RaftGRPCClient {

	grpcClient := RaftGRPCClient{}
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRaftClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutIn)
	grpcClient.client = c
	grpcClient.ctx = ctx
	defer cancel()

	return grpcClient
}

func (c RaftGRPCClient) ReplicateLog(lr *pb.LogRequest) (*pb.LogResponse, error) {
	return c.client.ReplicateLog(c.ctx, lr)
}

func (c RaftGRPCClient) RequestVoteFromPeer(vr *pb.VoteRequest) (*pb.VoteResponse, error) {
	return c.client.RequestVoteFromPeer(c.ctx, vr)
}

func (c RaftGRPCClient) BroadcastMessage(m *structpb.Struct) (*pb.BroadcastMessageResponse, error) {
	return c.client.BroadcastMessage(c.ctx, m)
}
