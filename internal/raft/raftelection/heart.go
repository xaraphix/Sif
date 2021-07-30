package raftelection

import "github.com/xaraphix/Sif/internal/raft"

type LeaderHeart struct {
}

func (l *LeaderHeart) StartBeating(rn *raft.RaftNode) {
	go func(n *raft.RaftNode) {
		for _, peer := range rn.Peers {
			rn.RPCAdapter.SendHeartbeatToPeer(peer)
		}
	}(rn)
}

func (l *LeaderHeart) StopBeating(rn *raft.RaftNode) {

}

func ifLeaderStartHeartbeatTransmitter(rn *raft.RaftNode) {
	if rn.CurrentRole == raft.LEADER {
		rn.Heart.StopBeating(rn)
	}
}
