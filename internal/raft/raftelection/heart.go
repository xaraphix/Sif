package raftelection

import (
	"time"

	"github.com/xaraphix/Sif/internal/raft"
)

var (
	LdrHrt raft.RaftHeart = &LeaderHeart{
		Heart: raft.Heart{
			DurationBetweenBeats: time.Millisecond * 200,
		},
	}
)

type LeaderHeart struct {
	raft.Heart
}

func (l *LeaderHeart) StartBeating(rn *raft.RaftNode) {
	rn.IsHeartBeating = true
	go func(n *raft.RaftNode) {
		for _, peer := range rn.Peers {
			rn.RPCAdapter.SendHeartbeatToPeer(peer)
		}
	}(rn)
}

func (l *LeaderHeart) StopBeating(rn *raft.RaftNode) {

}

func (l *LeaderHeart) IsBeating(rn *raft.RaftNode) bool {
	return rn.IsHeartBeating
}

func ifLeaderStartHeartbeatTransmitter(rn *raft.RaftNode) {
	if rn.CurrentRole == raft.LEADER {
		rn.Heart.StartBeating(rn)
	}
}
