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
}

type ElectionMonitor struct {
	ElectionTimeout       time.Duration
	LastELectionStartedAt time.Time
}

func (monitor *LeaderHeartbeatMonitor) Start() {

	go func(lhm *LeaderHeartbeatMonitor) {

		for {
			time.Sleep(monitor.HeartbeatTimeout)

			if time.Since(lhm.LeaderLastHeartbeat) >= lhm.HeartbeatTimeout {
				raftnode.Mu.Lock()
				raftnode.CurrentRole = CANDIDATE
				raftnode.Mu.Unlock()
			}

		}
	}(monitor)
}

func (em *ElectionMonitor) Start() {

}

func NewRaftNode() (*RaftNode, error) {

	once.Do(func() {
		Init()
	})

	return raftnode, nil
}

func Init() {
	raftnode = &RaftNode{}
	initializeRaftNode(raftnode)
	raftnode.LeaderHeartbeatMonitor.Start()
}

func initializeRaftNode(rn *RaftNode) {
	rn.CurrentRole = FOLLOWER
	rn.LeaderHeartbeatMonitor = &LeaderHeartbeatMonitor{
		LeaderLastHeartbeat: time.Time{},
		HeartbeatTimeout:    time.Duration(rand.Intn(150)+150) * time.Millisecond,
	}
	rn.ElectionMonitor = &ElectionMonitor{
		ElectionTimeout:       time.Duration(rand.Intn(150)+150) * time.Millisecond,
		LastELectionStartedAt: time.Time{},
	}
}
