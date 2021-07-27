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
	once        sync.Once
	raftnode    *RaftNode
	raftAdapter *RaftAdapter
)

type Node struct {
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
	Peers         []Peer
}

type RaftNode struct {
	Node
	ElectionInProgress     bool
	LeaderHeartbeatMonitor *LeaderHeartbeatMonitor
	ElectionMonitor        *ElectionMonitor
	Heart                  *Heart
}

type Heart struct {
	DurationBetweenBeats time.Duration
}

func (rn *RaftNode) StartLeaderHeartbeatMonitor() {
	turnOnLeaderHeartbeatMonitor(rn)
}

func (rn *RaftNode) StartElectionMonitor() {
	turnOnElectionMonitor(rn)
}

func (rn *RaftNode) StartElection() {
	startElection(rn)
}

func (rn *RaftNode) RequestVotes() {
	voteResponseChannel := make(chan VoteResponse, len(rn.Peers))
	requestVotesFromPeers(rn, voteResponseChannel)
	counter := 1

	for response := range voteResponseChannel {
		rn.UpdateVotesRecieved(response)
		if counter == len(rn.Peers) {
			break
		}
		counter = counter + 1
	}

	close(voteResponseChannel)

}

func (rn *RaftNode) GenerateVoteRequest() VoteRequest {
	return VoteRequest{}
}

func (rn *RaftNode) UpdateVotesRecieved(voteResponse VoteResponse) {
	//TODO
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
	rn.Peers = []Peer{{Id: 2, Address: ""}, {Id: 3, Address: ""}}
	rn.ElectionInProgress = false
	rn.LeaderHeartbeatMonitor = &LeaderHeartbeatMonitor{
		Monitor: Monitor{
			Stopped:         false,
			LastResetAt:     time.Time{},
			TimeoutDuration: time.Duration(rand.Intn(150)+150) * time.Millisecond,
		},
	}
	rn.ElectionMonitor = &ElectionMonitor{
		Monitor: Monitor{
			Stopped:         true,
			LastResetAt:     time.Time{},
			TimeoutDuration: time.Duration(rand.Intn(150)+150) * time.Millisecond,
		},
	}
}
