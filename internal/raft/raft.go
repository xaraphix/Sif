package raft

import (
	"time"
)

var (
	Sif *RaftNode
)

type Node struct {
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
	raftSignal         chan int

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
	Start(raftNode *RaftNode)
	Stop()
	Sleep()
	GetLastResetAt() time.Time
	Reset()
}

//go:generate mockgen -destination=mocks/mock_raftelection.go -package=mocks . RaftElection
type RaftElection interface {
	GetReceivedVotes() []VoteResponse
	StartElection(*RaftNode)
	GetResponseForVoteRequest(raftnode *RaftNode, voteRequest VoteRequest) VoteResponse
	GenerateVoteRequest(*RaftNode) VoteRequest
	GetLeaderHeartChannel() chan RaftNode
}

//go:generate mockgen -destination=mocks/mock_raftrpcadapter.go -package=mocks . RaftRPCAdapter
type RaftRPCAdapter interface {
	RequestVoteFromPeer(peer Peer, voteRequest VoteRequest) VoteResponse
	ReplicateLog(peer Peer, logRequest LogRequest)
	ReceiveLogRequest(logRequest LogRequest) LogResponse
	GenerateVoteResponse(VoteRequest) VoteResponse
	BroadcastMessage(msg map[string]interface{})
}

//go:generate mockgen -destination=mocks/mock_raftheart.go -package=mocks . RaftHeart
type RaftHeart interface {
	StopBeating(*RaftNode)
	StartBeating(*RaftNode)
	Sleep(*RaftNode)
}

//go:generate mockgen -destination=mocks/mock_raftlog.go -package=mocks . RaftLog
type RaftLog interface {
	GetLogs() []Log
	GetLog(rn *RaftNode, idx int32) Log
	ReplicateLog(raftNode *RaftNode, peer Peer)
	RespondToBroadcastMsgRequest(raftNode *RaftNode, msg map[string]interface{})
	RespondToLogRequest(raftNode *RaftNode, logRequest LogRequest) LogResponse
}

type RaftOptions struct {
	StartLeaderHeartbeatMonitorAfterInitializing bool
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

type VoteRequest struct {
	NodeId      int32
	CurrentTerm int32
	LogLength   int32
	LastTerm    int32
}

type VoteResponse struct {
	PeerId      int32
	Term        int32
	VoteGranted bool
}

type ElectionUpdates struct {
	ElectionOvertimed    bool
	ElectionCompleted    bool
	ElectionStopped      bool
	ElectionStarted      bool
	ElectionTimerStarted bool
}

type Peer struct {
	Id      int32  `yaml:"id"`
	Address string `yaml:"address"`
}

type LogRequest struct {
	LeaderId     int32
	Term  int32
	LogLength   int32
	LogTerm  int32
	CommitLength int32
	Entries      *[]Log
}

type LogResponse struct {
	FollowerId int32
	Term       int32
	AckLength  int32
	Success    bool
}

func NewRaftNode(
	fm RaftFile,
	rc RaftConfig,
	re RaftElection,
	lhm RaftMonitor,
	ra RaftRPCAdapter,
	l RaftLog,
	h RaftHeart,
	o RaftOptions,
) *RaftNode {
	DestructRaftNode(Sif)
	Sif = &RaftNode{
		FileMgr:                fm,
		Config:                 rc,
		ElectionMgr:            re,
		LeaderHeartbeatMonitor: lhm,
		RPCAdapter:             ra,
		LogMgr:                 l,
		Heart:                  h,
		raftSignal:             make(chan int),
	}

	raftnode := Sif
	initializeRaftNode(raftnode)
	if o.StartLeaderHeartbeatMonitorAfterInitializing {
		raftnode.LeaderHeartbeatMonitor.Start(raftnode)
	}
	return raftnode
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
	rn.SentLength = map[int32]int32{}
	rn.AckedLength = map[int32]int32{}
	rn.VotesReceived = nil
	rn.ElectionInProgress = false
}

func (rn *RaftNode) GetRaftSignalsChan() <-chan int {
	return rn.raftSignal
}

func (rn *RaftNode) SendSignal(signal int) {
	go func() {
		rn.raftSignal <- signal
	}()
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
