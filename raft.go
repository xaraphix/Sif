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

	ElectionTimeout       time.Duration
	HeartbeatTimeout      time.Duration
	LeaderHeartbeatRate   time.Duration
	LeaderLastHeartbeat   time.Time
	LastELectionStartedAt time.Time
}

func CreateRaftNode() (*RaftNode, error) {
	once.Do(func() {
		raftnode = &RaftNode{}
		initializeRaftNode(raftnode)
		startLeaderHeartbeatMonitor(raftnode)
	})

	return raftnode, nil
}

func initializeRaftNode(rn *RaftNode) {
	rn.CurrentRole = FOLLOWER
	rn.LeaderLastHeartbeat = time.Now()
	rn.HeartbeatTimeout = time.Duration(rand.Intn(150)+150) * time.Millisecond
}

func startLeaderHeartbeatMonitor(rn *RaftNode) {
	go func(n *RaftNode) {
		for {
			if time.Since(n.LeaderLastHeartbeat) >= n.HeartbeatTimeout {
				n.CurrentRole = CANDIDATE
			}

			time.Sleep(rn.HeartbeatTimeout)
		}
	}(rn)
}
