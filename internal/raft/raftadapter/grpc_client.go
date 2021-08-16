package raftadapter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type RaftGRPCClient struct {
	client          pb.RaftClient
	timeoutDuration time.Duration
}

func NewRaftGRPCClient(address string, timeoutIn time.Duration) *RaftGRPCClient {

	grpcClient := RaftGRPCClient{}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewRaftClient(conn)
	grpcClient.client = c
	grpcClient.timeoutDuration = timeoutIn*5

	return &grpcClient
}

func (c *RaftGRPCClient) ReplicateLog(lr *pb.LogRequest) (*pb.LogResponse, error) {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(c.timeoutDuration))
	x, y := c.client.ReplicateLog(ctx, lr)

	if x == nil {
		return nil, y
	}

	logrus.WithFields(logrus.Fields{
		"MyId": lr.LeaderId,
		"From": x.FollowerId,
	}).Info("Received log response")

	if y != nil {
		logrus.Error("In repl log :: " + y.Error())
		fmt.Printf(y.Error())
	}

	return x, nil
}

func (c *RaftGRPCClient) RequestVoteFromPeer(vr *pb.VoteRequest) (*pb.VoteResponse, error) {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(c.timeoutDuration))
	x, y := c.client.RequestVoteFromPeer(ctx, vr)
	if x == nil {
		return nil, y
	}

	logrus.WithFields(logrus.Fields{
		"Requested From": vr.NodeId,
		"Response by":    x.PeerId,
	}).Info("Received vote response")

	if y != nil {
		logrus.Error("In req vote :: " + y.Error())
		fmt.Printf(y.Error())
	}

	return x, nil
}

func (c *RaftGRPCClient) BroadcastMessage(m *structpb.Struct) (*pb.BroadcastMessageResponse, error) {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(c.timeoutDuration))
	x, y := c.client.BroadcastMessage(ctx, m)
	if y != nil {
		logrus.Error("In broadcast msg :: " + y.Error())
		fmt.Printf(y.Error())
	}

	return x, nil
}
