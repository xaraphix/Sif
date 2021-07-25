package raft

import (
	"math/rand"
	"sync"
	"time"
)

const (
	FOLLOWER  = "follower"
	CANDIDATE = "candidate"
)

var (
	once     sync.Once
	raftnode *RaftNode
)

type RaftNode struct {
	Mu sync.Mutex

	Id            int32
	CurrentTerm   int32
	CurrentRole   string
	CurrentLeader int32
	VotedFor      int32
	VotesReceived []int32
	CommitIndex   int32
	CommitLength  int32
	AckedLength   map[int32]int32
	SentLength    map[int32]int32
	LastApplied   int32
	NextIndex     int32
	MatchIndex    int32
	PrevLogIndex  int32
	// Peers         []Peer
	// Log           []*pb.Log

	LeaderHeartbeatMonitor *LeaderHeartbeatMonitor
	ElectionMonitor        *ElectionMonitor
	Heart                  *Heart
}

type Heart struct {
	DurationBetweenBeats time.Duration
}

type LeaderHeartbeatMonitor struct {
	HeartbeatTimeout    time.Duration
	LeaderLastHeartbeat time.Time
	Stopped             bool
}

type ElectionMonitor struct {
	ElectionTimeout       time.Duration
	LastELectionStartedAt time.Time
}

func (rn *RaftNode) StartLeaderHeartbeatMonitor() {

	go func(r *RaftNode) {

		for {

			if time.Since(r.LeaderHeartbeatMonitor.LeaderLastHeartbeat) >= r.LeaderHeartbeatMonitor.HeartbeatTimeout {
				r.Mu.Lock()
				r.CurrentRole = CANDIDATE
				r.VotedFor = raftnode.Id
				r.CurrentTerm = raftnode.CurrentTerm + 1
				r.Mu.Unlock()
			}

			if r.LeaderHeartbeatMonitor.Stopped {
				break
			}

			time.Sleep(r.LeaderHeartbeatMonitor.HeartbeatTimeout)
		}
	}(rn)
}

func (em *ElectionMonitor) Start() {

}

func NewRaftNode(returnExistingIfPresent bool) *RaftNode {
	if returnExistingIfPresent {
		once.Do(func() {
			raftnode = &RaftNode{}
			initializeRaftNode(raftnode)
			raftnode.StartLeaderHeartbeatMonitor()
		})
		return raftnode
	} else {
		raftnode = &RaftNode{}
		initializeRaftNode(raftnode)
		raftnode.StartLeaderHeartbeatMonitor()
		return raftnode
	}
}

func initializeRaftNode(rn *RaftNode) {
	rn.CurrentRole = FOLLOWER
	rn.CurrentTerm = 0
	rn.LeaderHeartbeatMonitor = &LeaderHeartbeatMonitor{
		Stopped:             false,
		LeaderLastHeartbeat: time.Time{},
		HeartbeatTimeout:    time.Duration(rand.Intn(150)+150) * time.Millisecond,
	}
	rn.ElectionMonitor = &ElectionMonitor{
		ElectionTimeout:       time.Duration(rand.Intn(150)+150) * time.Millisecond,
		LastELectionStartedAt: time.Time{},
	}
}
