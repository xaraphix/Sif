package raft

import (
	"sync"
	"time"

	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"google.golang.org/protobuf/types/known/structpb"
)

func init() {
	setupLogging()
}

var (
	Sif *RaftNode
)

type Node struct {
	LogAckMu sync.Mutex

	Id string

  //to be stored in non volatile storage
	CurrentTerm  int32
	CommitLength int32
	VotedFor     string
	Logs         []*pb.Log

  //can live in volatile storage
	CurrentRole   string
	CurrentLeader string
	VotesReceived []pb.VoteResponse
	CommitIndex   int32
	AckedLength   map[string]int32
	SentLength    map[string]int32
	LastApplied   int32
	NextIndex     int32
	MatchIndex    int32
	PrevLogIndex  int32
	Peers         []Peer
}

type RaftNode struct {
	Node

	ElectionInProgress bool
	StartedRPCAdapter  bool
	raftSignal         chan int
	HeartDone          chan bool
	ElectionDone       chan bool

	FileMgr                RaftFile
	Config                 RaftConfig
	ElectionMgr            RaftElection
	LeaderHeartbeatMonitor RaftMonitor
	RPCAdapter             RaftRPCAdapter
	LogMgr                 RaftLog
	Heart                  RaftHeart
}

//go:generate mockgen -destination=../../test/mocks/mock_raftfile.go -package=mocks . RaftFile
type RaftFile interface {
	LoadFile(filepath string) ([]byte, error)
	SaveFile(filepath string) error
}

//go:generate mockgen -destination=../../test/mocks/mock_raftconfig.go -package=mocks . RaftConfig
type RaftConfig interface {
	LoadConfig(*RaftNode)
	SetConfigFilePath(string)
	DidNodeCrash(*RaftNode) bool
	InstanceName() string
	InstanceId() string
	Peers() []Peer
	InstanceDirPath() string
	Version() string
	LogFilePath() string
	Logs() []*pb.Log
	CurrentTerm() int32
	CommitLength() int32
	VotedFor() string
	Host() string
	Port() string
}

//go:generate mockgen -destination=../../test/mocks/mock_raftmonitor.go -package=mocks . RaftMonitor
type RaftMonitor interface {
	Start(raftNode *RaftNode)
	Stop()
	Sleep()
	GetLastResetAt() time.Time
	Reset()
}

//go:generate mockgen -destination=../../test/mocks/mock_raftelection.go -package=mocks . RaftElection
type RaftElection interface {
	GetReceivedVotes() []*pb.VoteResponse
	StartElection(*RaftNode)
	ManageElection(*RaftNode) bool
	GetResponseForVoteRequest(raftnode *RaftNode, voteRequest *pb.VoteRequest) (*pb.VoteResponse, error)
	BecomeACandidate(*RaftNode)
	GenerateVoteRequest(*RaftNode) *pb.VoteRequest
	GetLeaderHeartChannel() chan *RaftNode
	GetElectionTimeoutDuration() time.Duration
	SetElectionTimeoutDuration(time.Duration)
}

//go:generate mockgen -destination=../../test/mocks/mock_raftrpcadapter.go -package=mocks . RaftRPCAdapter
type RaftRPCAdapter interface {
	RequestVoteFromPeer(Peer, *pb.VoteRequest) *pb.VoteResponse
	ReplicateLog(Peer, *pb.LogRequest) *pb.LogResponse
	BroadcastMessage(leader Peer, msg *structpb.Struct) *pb.BroadcastMessageResponse
	GetResponseForVoteRequest(raftnode *RaftNode, voteRequest *pb.VoteRequest) (*pb.VoteResponse, error)

	StartAdapter(*RaftNode)
	StopAdapter()
	GetRaftInfo(Peer, *pb.RaftInfoRequest) *pb.RaftInfoResponse
}

//go:generate mockgen -destination=../../test/mocks/mock_raftheart.go -package=mocks . RaftHeart
type RaftHeart interface {
	SetLeaderHeartbeatTimeout(time.Duration)
	StopBeating(*RaftNode)
	StartBeating(*RaftNode)
	Sleep(*RaftNode)
}

//go:generate mockgen -destination=../../test/mocks/mock_raftlog.go -package=mocks . RaftLog
type RaftLog interface {
	GetLogs() []*pb.Log
	GetLog(rn *RaftNode, idx int32) *pb.Log
	ReplicateLog(raftNode *RaftNode, peer Peer)
	RespondToBroadcastMsgRequest(raftNode *RaftNode, msg *structpb.Struct) (*pb.BroadcastMessageResponse, error)
	RespondToLogReplicationRequest(raftNode *RaftNode, logRequest *pb.LogRequest) (*pb.LogResponse, error)
}

type RaftOptions struct {
	StartLeaderHeartbeatMonitorAfterInitializing bool
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

type BroadcastMessageResponse struct {
}
type ElectionUpdates struct {
	ElectionOvertimed    bool
	ElectionCompleted    bool
	ElectionStopped      bool
	ElectionStarted      bool
	ElectionTimerStarted bool
}

type Peer struct {
	Id      string  `yaml:"id"`
	Address string `yaml:"address"`
}

type RaftDeps struct {
	FileManager      RaftFile
	ConfigManager    RaftConfig
	ElectionManager  RaftElection
	HeartbeatMonitor RaftMonitor
	RPCAdapter       RaftRPCAdapter
	LogManager       RaftLog
	Heart            RaftHeart
	Options          RaftOptions
}

func NewRaftNode(deps RaftDeps) *RaftNode {
	sif := &RaftNode{
		FileMgr:                deps.FileManager,
		Config:                 deps.ConfigManager,
		ElectionMgr:            deps.ElectionManager,
		LeaderHeartbeatMonitor: deps.HeartbeatMonitor,
		RPCAdapter:             deps.RPCAdapter,
		LogMgr:                 deps.LogManager,
		Heart:                  deps.Heart,
		raftSignal:             make(chan int),
	}

	raftnode := sif
	initializeRaftNode(raftnode)
	if deps.Options.StartLeaderHeartbeatMonitorAfterInitializing {
		raftnode.LeaderHeartbeatMonitor.Start(raftnode)
	}
	return raftnode
}

func (rn *RaftNode) Close() {
	rn.RPCAdapter.StopAdapter()
	if rn.CurrentRole == LEADER {
		rn.HeartDone <- true
	}

	if rn.ElectionInProgress {
		rn.ElectionDone <- true
	}

	rn.LeaderHeartbeatMonitor.Stop()
	close(rn.HeartDone)
	close(rn.ElectionDone)
	rn = nil
}

func initializeRaftNode(rn *RaftNode) {
	rn.Config.LoadConfig(rn)
	rn.Id = getId(rn)
	rn.CurrentRole = getCurrentRole(rn)
	rn.CurrentTerm = getCurrentTerm(rn)
	rn.Logs = getLogs(rn)
	rn.VotedFor = getVotedFor(rn)
	rn.CommitLength = getCommitLength(rn)
	rn.Peers = rn.Config.Peers()
	rn.SentLength = map[string]int32{}
	rn.AckedLength = map[string]int32{}
	rn.VotesReceived = nil
	rn.ElectionInProgress = false
	rn.RPCAdapter.StartAdapter(rn)
	rn.HeartDone = make(chan bool)
	rn.ElectionDone = make(chan bool)
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

func getId(rn *RaftNode) string {
	return rn.Config.InstanceId()
}

func getCurrentTerm(rn *RaftNode) int32 {
	if rn.Config.DidNodeCrash(rn) {
		return rn.Config.CurrentTerm()
	} else {
		return 0
	}
}

func getLogs(rn *RaftNode) []*pb.Log {
	if rn.Config.DidNodeCrash(rn) {
		return rn.Config.Logs()
	} else {
		logs := make([]*pb.Log, 0)
		return logs
	}
}

func getVotedFor(rn *RaftNode) string {
	if rn.Config.DidNodeCrash(rn) {
		return rn.Config.VotedFor()
	} else {
		return ""
	}
}

func getCommitLength(rn *RaftNode) int32 {
	if rn.Config.DidNodeCrash(rn) {
		return rn.Config.CommitLength()
	} else {
		return 0
	}
}

func (rn *RaftNode) GetPeerById(id string) Peer {
	for _, peer := range rn.Peers {
		if peer.Id == id {
			return peer
		}
	}
	return Peer{}
}
