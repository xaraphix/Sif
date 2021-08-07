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

	go func(n *raft.RaftNode) {
		for _, peer := range rn.Peers {
			rn.SentLength[peer.Id] = int32(len(rn.Logs))
			rn.AckedLength[peer.Id] = 0
			rn.LogMgr.ReplicateLog(rn, peer)
		}
	}(rn)
	go startBeating(rn)
}

func startBeating(rn *raft.RaftNode) {
	rn.IsHeartBeating = true
	for {
		if rn.IsHeartBeating == false && rn.CurrentRole == raft.LEADER {
			break
		} else {
			go func(n *raft.RaftNode) {
				for _, peer := range rn.Peers {
					rn.LogMgr.ReplicateLog(rn, peer)
				}
			}(rn)

			rn.Heart.Sleep(rn)
		}
	}
}

func (l *LeaderHeart) Sleep(rn *raft.RaftNode) {
	time.Sleep(l.DurationBetweenBeats)
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
