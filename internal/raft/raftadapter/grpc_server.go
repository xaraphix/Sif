package raftadapter

import (
	"context"
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/xaraphix/Sif/internal/raft"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCServer struct {
	pb.UnimplementedRaftServer
	raftnode *raft.RaftNode
	server   *grpc.Server
}

func NewGRPCServer(rn *raft.RaftNode) *GRPCServer {
	s := &GRPCServer{}
	s.raftnode = rn
	return s
}

func (s *GRPCServer) Start(host string, port string) {
	go func() {
		address := host + ":" + port
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		logrus.WithFields(logrus.Fields{
			"Address": address,
		}).Debug("Starting GRPC Server")

		s.server = grpc.NewServer()
		pb.RegisterRaftServer(s.server, s)
		s.server.Serve(lis)
	}()
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

func (s *GRPCServer) RequestVoteFromPeer(ctx context.Context, vr *pb.VoteRequest) (*pb.VoteResponse, error) {

	logrus.WithFields(logrus.Fields{
		"From":        vr.NodeId,
		"Received by": s.raftnode.Id,
	}).Debug("Received Vote Request")

	return s.raftnode.RPCAdapter.GetResponseForVoteRequest(s.raftnode, vr)
}

func (s *GRPCServer) ReplicateLog(ctx context.Context, vr *pb.LogRequest) (*pb.LogResponse, error) {

	logrus.WithFields(logrus.Fields{
		"From":        vr.LeaderId,
		"Received by": s.raftnode.Id,
	}).Debug("Replicate Log Receieved")

	return s.raftnode.LogMgr.RespondToLogReplicationRequest(s.raftnode, vr)
}

func (s *GRPCServer) BroadcastMessage(ctx context.Context, msg *structpb.Struct) (*pb.BroadcastMessageResponse, error) {
	return s.raftnode.LogMgr.RespondToBroadcastMsgRequest(s.raftnode, msg)
}

func (s *GRPCServer) GetRaftInfo(ctx context.Context, in *pb.RaftInfoRequest) (*pb.RaftInfoResponse, error) {
	return &pb.RaftInfoResponse{
		NodeId:        s.raftnode.Id,
		CurrentLeader: s.raftnode.CurrentLeader,
		CurrentTerm:   s.raftnode.CurrentTerm,
	}, nil
}
