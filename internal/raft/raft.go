package raft

import (
	"sync"
	"time"
)

const (
	FOLLOWER  = "follower"
	CANDIDATE = "candidate"
	LEADER    = "leader"
)

var (
	once     sync.Once
	raftnode *RaftNode
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
	ElectionInProgress bool

	Config                 RaftConfig
	ElectionManager        RaftElection
	LeaderHeartbeatMonitor RaftMonitor
	ElectionMonitor        RaftMonitor
	RPCAdapter             RaftRPCAdapter
	Heart                  Heart
}

type Heart struct {
	DurationBetweenBeats time.Duration
}

type Monitor struct {
	TimeoutDuration time.Duration
	LastResetAt     time.Time
	Stopped         bool
	Started         bool
	AutoStart       bool
}

type Peer struct {
	Id      int32  `yaml:"id"`
	Address string `yaml:"address"`
}

//go:generate mockgen -destination=mocks/mock_raftconfig.go -package=mocks . RaftConfig
type RaftConfig interface {
	LoadConfig() RaftConfig
	DidNodeCrash() bool
	InstanceName() string
	InstanceId() int32
	Peers() []Peer
	InstanceDirPath() string
	Version() string
}

//go:generate mockgen -destination=mocks/mock_raftmonitor.go -package=mocks . RaftMonitor
type RaftMonitor interface {
	Start(*RaftNode)
	IsAutoStartOn() bool
	Stop()
	Sleep()
	GetLastResetAt() time.Time
}

//go:generate mockgen -destination=mocks/mock_raftelection.go -package=mocks . RaftElection
type RaftElection interface {
	StartElection(*RaftNode)
	RequestVotes(*RaftNode)
	StopElection(*RaftNode)
	GenerateVoteRequest(*RaftNode) VoteRequest
}

//go:generate mockgen -destination=mocks/mock_raftrpcadapter.go -package=mocks . RaftRPCAdapter
type RaftRPCAdapter interface {
	RequestVoteFromPeer(peer Peer, voteRequest VoteRequest) VoteResponse
	SendHeartbeatToPeer()
}

func NewRaftNode(
	rc RaftConfig,
	re RaftElection,
	lhm *LeaderHeartbeatMonitor,
	em *ElectionMonitor,
	ra RaftRPCAdapter,
	forceNew bool,
) *RaftNode {
	rn := &RaftNode{
		Config:                 rc,
		ElectionManager:        re,
		LeaderHeartbeatMonitor: lhm,
		ElectionMonitor:        em,
		RPCAdapter:             ra,
	}

	if forceNew {
		raftnode = rn
		initializeRaftNode(raftnode)
		if raftnode.LeaderHeartbeatMonitor.IsAutoStartOn() {
			raftnode.LeaderHeartbeatMonitor.Start(raftnode)
		}
		return raftnode
	} else {

		if raftnode == nil {
			raftnode = rn
			initializeRaftNode(raftnode)
			raftnode.LeaderHeartbeatMonitor.Start(raftnode)
		}
		return raftnode
	}
}

func DestructRaftNode(rn *RaftNode) {
	rn = nil
}

func initializeRaftNode(rn *RaftNode) {
	rn.Config = rn.Config.LoadConfig()
	rn.CurrentRole = getCurrentRole(rn)
	rn.CurrentTerm = getCurrentTerm(rn)
	rn.Peers = rn.Config.Peers()
	rn.VotesReceived = nil
	rn.ElectionInProgress = false
}

func (m *Monitor) Stop() {
	m.Stopped = true
}

func (m *Monitor) GetLastResetAt() time.Time {
	return m.LastResetAt
}

func (m *Monitor) Sleep() {
	time.Sleep(m.TimeoutDuration)
}

func (m *Monitor) IsAutoStartOn() bool {
	return m.AutoStart
}

type VoteRequest struct {
}

type VoteResponse struct {
	PeerId      int32
	VoteGranted bool
}

func getCurrentRole(rn *RaftNode) string {
	//TODO
	// check if node crashed
	// 			check if it was a leader

	return FOLLOWER
}

func getCurrentTerm(rn *RaftNode) int32 {
	//TODO
	// check if node crashed
	// 			check if it was a leader

	return 0
}