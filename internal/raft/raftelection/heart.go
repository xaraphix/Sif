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
  timer := time.NewTicker(l.DurationBetweenBeats)
	go func(n *raft.RaftNode) {
		n.LogEvent(raft.HeartbeatStarted, raft.RaftEventDetails{CurrentTerm: n.CurrentTerm, CurrentRole: n.CurrentRole})
		for {
			select {
			case <-n.HeartDone:
        timer.Stop()
				break
			case <- timer.C:
				if n.CurrentRole != raft.LEADER {
					n.LogEvent(raft.HeartbeatStopped, raft.RaftEventDetails{CurrentTerm: n.CurrentTerm, CurrentRole: n.CurrentRole, CurrentLeader: n.CurrentLeader})
          timer.Stop()
					break
				} else {
					go func(node *raft.RaftNode) {
						for _, peer := range node.Peers {
							node.LogMgr.ReplicateLog(node, peer)
						}
					}(n)
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
