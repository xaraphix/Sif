package raftadapter

import (
	"github.com/xaraphix/Sif/internal/raft"
	. "github.com/xaraphix/Sif/internal/raft"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"google.golang.org/protobuf/types/known/structpb"
)

type RaftNodeAdapter struct {
	clients map[int32]*RaftGRPCClient
}

func NewRaftNodeAdapter() *RaftNodeAdapter {
	return &RaftNodeAdapter{}
}

func (a *RaftNodeAdapter) StartAdapter(rn *raft.RaftNode) {
	a.initializeClients(rn)
	grpcServer := NewGRPCServer(rn)
	grpcServer.Start(rn.Config.Host(), rn.Config.Port())
}

func (a *RaftNodeAdapter) initializeClients(rn *raft.RaftNode) {
	a.clients = make(map[int32]*RaftGRPCClient)
	for _, peer := range rn.Peers {
		a.clients[peer.Id] =NewRaftGRPCClient(peer.Address, rn.ElectionMgr.GetElectionTimeoutDuration())
	}
}

func (a *RaftNodeAdapter) RequestVoteFromPeer(peer Peer, voteRequest *pb.VoteRequest) *pb.VoteResponse {
	r, _ := a.clients[peer.Id].RequestVoteFromPeer(voteRequest)
	return r
}

func (a *RaftNodeAdapter) ReplicateLog(peer Peer, logRequest *pb.LogRequest) *pb.LogResponse {
	r, _ := a.clients[peer.Id].ReplicateLog(logRequest)
	return r
}

func (a *RaftNodeAdapter) BroadcastMessage(peer Peer, msg *structpb.Struct) *pb.BroadcastMessageResponse {
	r, _ := a.clients[peer.Id].BroadcastMessage(msg)
	return r
}
