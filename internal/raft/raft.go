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
	VotesReceived []VoteResponse
	CommitIndex   int32
	CommitLength  int32
	AckedLength   map[int32]int32
	SentLength    map[int32]int32
	LastApplied   int32
	NextIndex     int32
	MatchIndex    int32
	PrevLogIndex  int32
	Peers         []Peer
	Logs          []Log
}

type RaftNode struct {
	Node

	ElectionInProgress bool
	IsHeartBeating     bool

	FileMgr                RaftFile
	Config                 RaftConfig
	ElectionMgr            RaftElection
	LeaderHeartbeatMonitor RaftMonitor
	RPCAdapter             RaftRPCAdapter
	LogMgr                 RaftLog
	Heart                  RaftHeart
}

//go:generate mockgen -destination=mocks/mock_raftfile.go -package=mocks . RaftFile
type RaftFile interface {
	LoadFile(filepath string) ([]byte, error)
	SaveFile(filepath string) error
}

//go:generate mockgen -destination=mocks/mock_raftconfig.go -package=mocks . RaftConfig
type RaftConfig interface {
	LoadConfig(*RaftNode)
	DidNodeCrash(*RaftNode) bool
	InstanceName() string
	InstanceId() int32
	Peers() []Peer
	InstanceDirPath() string
	Version() string
	LogFilePath() string
	Logs() *[]Log
	CurrentTerm() int32
	CommitLength() int32
	VotedFor() int32
}

//go:generate mockgen -destination=mocks/mock_raftmonitor.go -package=mocks . RaftMonitor
type RaftMonitor interface {
	Start(raftNode *RaftNode, electionChannel chan ElectionUpdates)
	Stop()
	Sleep()
	GetLastResetAt() time.Time
}

//go:generate mockgen -destination=mocks/mock_raftelection.go -package=mocks . RaftElection
type RaftElection interface {
	StartElection(*RaftNode)
	RequestVotes(raftnode *RaftNode, electionChannel chan ElectionUpdates)
	StopElection(*RaftNode)
	GenerateVoteRequest(*RaftNode) VoteRequest
}

//go:generate mockgen -destination=mocks/mock_raftrpcadapter.go -package=mocks . RaftRPCAdapter
type RaftRPCAdapter interface {
	RequestVoteFromPeer(peer Peer, voteRequest VoteRequest) VoteResponse
	SendHeartbeatToPeer(peer Peer)
}

//go:generate mockgen -destination=mocks/mock_raftheart.go -package=mocks . RaftHeart
type RaftHeart interface {
	StopBeating(*RaftNode)
	StartBeating(*RaftNode)
	IsBeating(*RaftNode) bool
}

//go:generate mockgen -destination=mocks/mock_raftlog.go -package=mocks . RaftLog
type RaftLog interface {
	GetLogs() []Log
	GetLog(idx int32) Log
}

type Log struct {
	Term    int32                  `json:"term"`
	Message map[string]interface{} `json:"message"`
}

type Heart struct {
	DurationBetweenBeats time.Duration
}

type Monitor struct {
	TimeoutDuration time.Duration
	LastResetAt     time.Time
	Stopped         bool
	Started         bool
}

type ElectionUpdates struct {
	ElectionOvertimed bool
}

type Peer struct {
	Id      int32  `yaml:"id"`
	Address string `yaml:"address"`
}

func NewRaftNode(
	fm RaftFile,
	rc RaftConfig,
	re RaftElection,
	lhm *LeaderHeartbeatMonitor,
	ra RaftRPCAdapter,
	l RaftLog,
	h RaftHeart,
	forceNew bool,
) *RaftNode {
	rn := &RaftNode{
		FileMgr:                fm,
		Config:                 rc,
		ElectionMgr:            re,
		LeaderHeartbeatMonitor: lhm,
		RPCAdapter:             ra,
		LogMgr:                 l,
		Heart:                  h,
	}

	if forceNew {
		raftnode = rn
		initializeRaftNode(raftnode)
		raftnode.LeaderHeartbeatMonitor.Start(raftnode, nil)
		return raftnode
	} else {

		if raftnode == nil {
			raftnode = rn
			initializeRaftNode(raftnode)
			raftnode.LeaderHeartbeatMonitor.Start(raftnode, nil)
		}
		return raftnode
	}
}

func DestructRaftNode(rn *RaftNode) {
	rn = nil
}

func initializeRaftNode(rn *RaftNode) {
	rn.Config.LoadConfig(rn)
	rn.CurrentRole = getCurrentRole(rn)
	rn.CurrentTerm = getCurrentTerm(rn)
	rn.Logs = *getLogs(rn)
	rn.VotedFor = getVotedFor(rn)
	rn.CommitLength = getCommitLength(rn)
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

type VoteRequest struct {
}

type VoteResponse struct {
	PeerId      int32
	VoteGranted bool
}

func getCurrentRole(rn *RaftNode) string {
	return FOLLOWER
}

func getCurrentTerm(rn *RaftNode) int32 {
	if rn.Config.DidNodeCrash(rn) {
		return rn.Config.CurrentTerm()
	} else {
		return 0
	}
}

func getLogs(rn *RaftNode) *[]Log {
	if rn.Config.DidNodeCrash(rn) {
		return rn.Config.Logs()
	} else {
		return &[]Log{}
	}
}

func getVotedFor(rn *RaftNode) int32 {
	if rn.Config.DidNodeCrash(rn) {
		return rn.Config.VotedFor()
	} else {
		return 0
	}
}

func getCommitLength(rn *RaftNode) int32 {
	if rn.Config.DidNodeCrash(rn) {
		return rn.Config.CommitLength()
	} else {
		return 0
	}
}
