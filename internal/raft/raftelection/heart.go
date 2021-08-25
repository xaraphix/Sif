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

func (l *LeaderHeart) SetLeaderHeartbeatTimeout(duration time.Duration) {
	l.DurationBetweenBeats = duration
}

func (l *LeaderHeart) StartBeating(rn *raft.RaftNode) {
	go func(rn *raft.RaftNode) {
		raft.LogEvent(raft.HeartbeatStarted, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
		for {
			select {
			case <-rn.HeartDone:
				break
			default:
				if rn.CurrentRole != raft.LEADER {
					raft.LogEvent(raft.HeartbeatStopped, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole, CurrentLeader: rn.CurrentLeader})
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
	}(rn)
}

func (l *LeaderHeart) Sleep(rn *raft.RaftNode) {
	time.Sleep(l.DurationBetweenBeats)
}

func (l *LeaderHeart) StopBeating(rn *raft.RaftNode) {

}
