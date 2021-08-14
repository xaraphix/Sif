package raftelection

import (
	"time"

	"github.com/xaraphix/Sif/internal/raft"
)

type LeaderHeart struct {
	raft.Heart
}

func NewLeaderHeart() *LeaderHeart {
	return &LeaderHeart{
		Heart: raft.Heart{
			DurationBetweenBeats: time.Millisecond * 200,
		},
	}
}

func (l *LeaderHeart) StartBeating(rn *raft.RaftNode) {
	go func(rn *raft.RaftNode) {
		rn.SendSignal(raft.HeartbeatStarted)
		for {
			if rn.CurrentRole != raft.LEADER {
				rn.SendSignal(raft.HeartbeatStopped)
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
	}(rn)
}

func (l *LeaderHeart) Sleep(rn *raft.RaftNode) {
	time.Sleep(l.DurationBetweenBeats)
}

func (l *LeaderHeart) StopBeating(rn *raft.RaftNode) {

}
