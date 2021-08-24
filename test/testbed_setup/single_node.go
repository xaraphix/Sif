package testbed_setup

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/xaraphix/Sif/internal/raft"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"github.com/xaraphix/Sif/internal/raft/raftconfig"
	"github.com/xaraphix/Sif/internal/raft/raftelection"
	"github.com/xaraphix/Sif/internal/raft/raftfile"
	"github.com/xaraphix/Sif/internal/raft/raftlog"
	"github.com/xaraphix/Sif/test/mocks"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v2"
)

func SetupRaftNode(preNodeSetupCB func(
	fileMgr *mocks.MockRaftFile,
	logMgr *mocks.MockRaftLog,
	election *mocks.MockRaftElection,
	adapter *mocks.MockRaftRPCAdapter,
	heart *mocks.MockRaftHeart,
	monitor *mocks.MockRaftMonitor,
), options SetupOptions) MockSetupVars {

	var (
		node   *raft.RaftNode
		term_0 int32

		election   raft.RaftElection
		config     raft.RaftConfig
		fileMgr    raft.RaftFile
		logMgr     raft.RaftLog
		rpcAdapter raft.RaftRPCAdapter
		heart      raft.RaftHeart
		monitor    raft.RaftMonitor
	)

	mockRPCAdapter, RpcCtrl := GetMockRPCAdapter()
	mockHeart, HeartCtrl := GetMockHeart()
	mockElection, ElectionCtrl := GetMockElection()
	mockFile, FileCtrl := GetMockFile()
	mockLog, LogCtrl := GetMockLog()
	mockMonitor, MonitorCtrl := GetMockMonitor()

	ctrls := Controllers{
		RpcCtrl:      RpcCtrl,
		LogCtrl:      LogCtrl,
		HeartCtrl:    HeartCtrl,
		FileCtrl:     FileCtrl,
		ElectionCtrl: ElectionCtrl,
		MonitorCtrl:  MonitorCtrl,
	}

	raftOptions := raft.RaftOptions{}

	if options.StartMonitor {
		raftOptions.StartLeaderHeartbeatMonitorAfterInitializing = true
	}

	election = raftelection.NewElectionManager()
	heart = raftelection.NewLeaderHeart()
	logMgr = &raftlog.LogMgr{}
	monitor = raft.NewLeaderHeartbeatMonitor(true)

	if options.MockFile {
		fileMgr = mockFile
	}

	if options.MockElection {
		election = mockElection
	}

	if options.MockHeart {
		heart = mockHeart
	}

	if options.MockLeaderMonitor {
		monitor = mockMonitor
	}

	if options.MockRPCAdapter {
		rpcAdapter = mockRPCAdapter
	}

	config = raftconfig.NewConfig()
	preNodeSetupCB(mockFile, mockLog, mockElection, mockRPCAdapter, mockHeart, mockMonitor)

	node = nil
	deps := raft.RaftDeps{
		FileManager:      fileMgr,
		ConfigManager:    config,
		ElectionManager:  election,
		HeartbeatMonitor: monitor,
		RPCAdapter:       rpcAdapter,
		LogManager:       logMgr,
		Heart:            heart,
		Options:          raftOptions,
	}

	node = raft.NewRaftNode(deps)
	term_0 = node.CurrentTerm

	return MockSetupVars{
		Node:   node,
		Ctrls:  ctrls,
		Term_0: term_0,
	}
}

func GetMockRPCAdapter() (*mocks.MockRaftRPCAdapter, *gomock.Controller) {

	var (
		mockCtrl   *gomock.Controller
		rpcAdapter *mocks.MockRaftRPCAdapter
	)

	mockCtrl = gomock.NewController(GinkgoT())
	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	return rpcAdapter, mockCtrl
}

func GetMockHeart() (*mocks.MockRaftHeart, *gomock.Controller) {

	var (
		mockCtrl  *gomock.Controller
		mockHeart *mocks.MockRaftHeart
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockHeart = mocks.NewMockRaftHeart(mockCtrl)
	return mockHeart, mockCtrl
}

func GetMockConfig() (*mocks.MockRaftConfig, *gomock.Controller) {

	var (
		mockCtrl   *gomock.Controller
		mockConfig *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockConfig = mocks.NewMockRaftConfig(mockCtrl)
	return mockConfig, mockCtrl
}

func GetMockFile() (*mocks.MockRaftFile, *gomock.Controller) {

	var (
		mockCtrl *gomock.Controller
		mockFile *mocks.MockRaftFile
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockFile = mocks.NewMockRaftFile(mockCtrl)
	return mockFile, mockCtrl
}

func GetMockLog() (*mocks.MockRaftLog, *gomock.Controller) {

	var (
		mockCtrl *gomock.Controller
		mockLog  *mocks.MockRaftLog
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockLog = mocks.NewMockRaftLog(mockCtrl)
	return mockLog, mockCtrl
}

func GetMockMonitor() (*mocks.MockRaftMonitor, *gomock.Controller) {

	var (
		mockCtrl *gomock.Controller
		mockLog  *mocks.MockRaftMonitor
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockLog = mocks.NewMockRaftMonitor(mockCtrl)
	return mockLog, mockCtrl
}

func GetMockElection() (*mocks.MockRaftElection, *gomock.Controller) {
	var (
		mockCtrl     *gomock.Controller
		mockElection *mocks.MockRaftElection
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockElection = mocks.NewMockRaftElection(mockCtrl)
	return mockElection, mockCtrl
}

type SetupOptions struct {
	MockHeart         bool
	MockLog           bool
	MockElection      bool
	MockFile          bool
	MockConfig        bool
	MockRPCAdapter    bool
	MockLeaderMonitor bool
	StartMonitor      bool
}

type Controllers struct {
	RpcCtrl      *gomock.Controller
	LogCtrl      *gomock.Controller
	ConfigCtrl   *gomock.Controller
	FileCtrl     *gomock.Controller
	HeartCtrl    *gomock.Controller
	ElectionCtrl *gomock.Controller
	MonitorCtrl  *gomock.Controller
}

type MockSetupVars struct {
	Node                  *raft.RaftNode
	Term_0                int32
	SentHeartbeats        *map[int]bool
	SentVoteRequests      *map[int]*pb.VoteRequest
	ReceivedVoteResponse  *map[int32]*pb.VoteResponse
	SentLogReplicationReq **pb.LogRequest
	ReceivedLogResponse   *map[int32]*pb.LogResponse
	Ctrls                 Controllers
	LeaderId              int32
}

func SetupRaftNodeBootsUpFromCrash() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockLog:        true,
		MockFile:       true,
		MockRPCAdapter: true,
		StartMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {
		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, nil)
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))

		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	return SetupRaftNode(preNodeSetupCB, options)
}

func SetupRaftNodeInitialization() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockFile:       true,
		MockLog:        true,
		MockRPCAdapter: true,
		StartMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {
		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()

		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	return SetupRaftNode(preNodeSetupCB, options)
}

func SetupLeaderHeartbeatTimeout() MockSetupVars {
	options := SetupOptions{
		MockFile:       true,
		MockHeart:      false,
		MockLog:        true,
		MockElection:   false,
		MockRPCAdapter: true,
		StartMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {
		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	setupVars := SetupRaftNode(preNodeSetupCB, options)
	for {
		if setupVars.Node.ElectionInProgress == true {
			setupVars.Node.LeaderHeartbeatMonitor.Stop()
			break
		}
	}

	return setupVars
}

func SetupCandidateRequestsVoteFromCandidate() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockFile:       true,
		MockLog:        true,
		MockRPCAdapter: true,
		StartMonitor:   false,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()

	}
	return SetupRaftNode(preNodeSetupCB, options)
}

func SetupFollowerReceivesLogReplicationRequest() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockFile:       true,
		MockLog:        true,
		MockRPCAdapter: true,
		StartMonitor:   false,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()

	}

	logReplicationReq := &pb.LogRequest{
		LeaderId:     999,
		CurrentTerm:  0,
		SentLength:   0,
		PrevLogTerm:  0,
		CommitLength: 2,
		Entries:      Get2LogEntries(),
	}

	setupVars := SetupRaftNode(preNodeSetupCB, options)
	setupVars.SentLogReplicationReq = &logReplicationReq
	return setupVars
}

func SetupLeaderReceivingLogReplicationAck() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockFile:       true,
		MockLog:        true,
		MockRPCAdapter: true,
		StartMonitor:   false,
	}

	logResponseMap := make(map[int32]*pb.LogResponse)

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		lrPeer1 := &pb.LogResponse{
			FollowerId: testConfig.Peers()[0].Id,
			Term:       1,
			AckLength:  2,
			Success:    true,
		}
		lrPeer2 := &pb.LogResponse{
			FollowerId: testConfig.Peers()[1].Id,
			Term:       1,
			AckLength:  2,
			Success:    true,
		}
		logResponseMap[testConfig.Peers()[0].Id] = lrPeer1
		logResponseMap[testConfig.Peers()[1].Id] = lrPeer2

		adapter.EXPECT().ReplicateLog(testConfig.RaftPeers[0], gomock.Any()).Return(logResponseMap[testConfig.RaftPeers[0].Id]).Times(1)
		adapter.EXPECT().ReplicateLog(testConfig.RaftPeers[1], gomock.Any()).Return(logResponseMap[testConfig.RaftPeers[1].Id]).Times(1)
	}

	setupVars := SetupRaftNode(preNodeSetupCB, options)
	setupVars.ReceivedLogResponse = &logResponseMap
	setupVars.Node.CurrentRole = raft.LEADER
	setupVars.Node.CurrentTerm = 1
	setupVars.Node.Logs = Get2LogEntries()
	return setupVars
}


func SetupLogReplicationNotAckByFollower() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockFile:       true,
		MockLog:        true,
		MockRPCAdapter: true,
		StartMonitor:   false,
	}

	logResponseMap := make(map[int32]*pb.LogResponse)

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		lrPeer1 := &pb.LogResponse{
			FollowerId: testConfig.Peers()[0].Id,
			Term:       1,
			AckLength:  0,
			Success:    false,
		}
		lrPeer2 := &pb.LogResponse{
			FollowerId: testConfig.Peers()[1].Id,
			Term:       1,
			AckLength:  0,
			Success:    false,
		}
		logResponseMap[testConfig.Peers()[0].Id] = lrPeer1
		logResponseMap[testConfig.Peers()[1].Id] = lrPeer2

		adapter.EXPECT().ReplicateLog(testConfig.RaftPeers[0], gomock.Any()).Return(logResponseMap[testConfig.RaftPeers[0].Id]).AnyTimes()
		adapter.EXPECT().ReplicateLog(testConfig.RaftPeers[1], gomock.Any()).Return(logResponseMap[testConfig.RaftPeers[1].Id]).AnyTimes()
	}

	setupVars := SetupRaftNode(preNodeSetupCB, options)
	setupVars.ReceivedLogResponse = &logResponseMap
	setupVars.Node.CurrentRole = raft.LEADER
	setupVars.Node.CurrentTerm = 1
	setupVars.Node.Logs = Get2LogEntries()
	return setupVars
}

func SetupMajorityVotesAgainst() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockFile:       true,
		MockLog:        true,
		MockRPCAdapter: true,
		StartMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).DoAndReturn(func(interface{}, interface{}) *pb.VoteResponse {

			return &pb.VoteResponse{
				VoteGranted: false,
				Term:        1,
				PeerId:      testConfig.Peers()[0].Id,
			}

		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).DoAndReturn(func(interface{}, interface{}) *pb.VoteResponse {

			return &pb.VoteResponse{
				VoteGranted: false,
				Term:        1,
				PeerId:      testConfig.Peers()[1].Id,
			}

		}).AnyTimes()
		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}
	return SetupRaftNode(preNodeSetupCB, options)
}

func SetupMajorityVotesInFavor() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockFile:       true,
		MockLog:        true,
		MockRPCAdapter: true,
		StartMonitor:   true,
	}

	var logReplicationReq *pb.LogRequest

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: true,
			Term:        1,
			PeerId:      testConfig.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[1].Id,
		}).AnyTimes()

		adapter.EXPECT().ReplicateLog(testConfig.Peers()[0], gomock.Any()).Do(func(p raft.Peer, r *pb.LogRequest) *pb.LogResponse {
			logReplicationReq = r
			return &pb.LogResponse{}
		}).AnyTimes()

		adapter.EXPECT().ReplicateLog(testConfig.Peers()[1], gomock.Any()).Do(func(interface{}, interface{}) {
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	setupVars := SetupRaftNode(preNodeSetupCB, options)
	setupVars.SentLogReplicationReq = &logReplicationReq
	return setupVars
}

func SetupLeaderSendsHeartbeatsOnElectionConclusion() MockSetupVars {
	sentHeartbeats := &map[int]bool{
		2: false,
		3: false,
	}

	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockFile:       true,
		MockLog:        true,
		MockRPCAdapter: true,
		StartMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: true,
			Term:        1,
			PeerId:      testConfig.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[1].Id,
		}).AnyTimes()

		adapter.EXPECT().ReplicateLog(testConfig.Peers()[0], gomock.Any()).Do(func(interface{}, interface{}) {
			(*sentHeartbeats)[int(testConfig.Peers()[0].Id)] = true
		}).MinTimes(1)

		adapter.EXPECT().ReplicateLog(gomock.Any(), gomock.Any()).Do(func(interface{}, interface{}) {
			(*sentHeartbeats)[int(testConfig.Peers()[1].Id)] = true
		}).MinTimes(1)

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	setupVars := SetupRaftNode(preNodeSetupCB, options)
	return MockSetupVars{
		Node:           setupVars.Node,
		Term_0:         setupVars.Term_0,
		Ctrls:          setupVars.Ctrls,
		SentHeartbeats: sentHeartbeats,
	}
}

func SetupPeerTakesTooMuchTimeToRespond() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockFile:       true,
		MockLog:        true,
		MockRPCAdapter: true,
		StartMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Do(
			func(p raft.Peer, vr *pb.VoteRequest) *pb.VoteResponse {
				time.Sleep(100 * time.Second)
				return &pb.VoteResponse{}
			}).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[1].Id,
		},
		).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	return SetupRaftNode(preNodeSetupCB, options)
}

func SetupGettingLeaderHeartbeatDuringElection() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockLog:        true,
		MockFile:       true,
		MockRPCAdapter: true,
		StartMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}
	return SetupRaftNode(preNodeSetupCB, options)
}

func SetupFindingOtherLeaderThroughVoteResponses() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockFile:       true,
		MockLog:        true,
		MockRPCAdapter: true,
		StartMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {
		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))

		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(&pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	return SetupRaftNode(preNodeSetupCB, options)
}

func SetupCandidateReceivesVoteResponseWithHigherTerm() MockSetupVars {
	options := SetupOptions{
		MockHeart:      false,
		MockElection:   false,
		MockLog:        true,
		MockFile:       true,
		MockRPCAdapter: true,
		StartMonitor:   true,
	}

	sentVoteResponse := make(map[int32]*pb.VoteResponse)
	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		vr1 := &pb.VoteResponse{
			VoteGranted: false,
			Term:        10,
			PeerId:      testConfig.Peers()[0].Id,
		}

		vr2 := &pb.VoteResponse{
			VoteGranted: false,
			Term:        1,
			PeerId:      testConfig.Peers()[1].Id,
		}
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Return(vr1).AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(vr2).AnyTimes()
		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()

		sentVoteResponse[testConfig.Peers()[0].Id] = vr1
		sentVoteResponse[testConfig.Peers()[1].Id] = vr2
	}

	setupVars := SetupRaftNode(preNodeSetupCB, options)
	setupVars.ReceivedVoteResponse = &sentVoteResponse
	return setupVars
}

func SetupPeerReceivingCandidateVoteRequest() MockSetupVars {

	options := SetupOptions{
		MockHeart:         false,
		MockElection:      false,
		MockLog:           false,
		MockFile:          true,
		MockRPCAdapter:    true,
		MockLeaderMonitor: false,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := LoadTestRaftConfig()
		testPersistentStorageFile, _ := LoadTestRaftPersistentStorageFile()
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().StopAdapter().Return().AnyTimes()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(LoadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
	}

	setupVars := SetupRaftNode(preNodeSetupCB, options)

	return setupVars
}

func GetConfig() raftconfig.Config {
	return raftconfig.Config{}
}

func LoadTestRaftConfigFile() ([]byte, error) {
	_, testFile, _, _ := runtime.Caller(0)
	dir, err1 := filepath.Abs(filepath.Dir(testFile))
	if err1 != nil {
		log.Fatal(err1)
	}
	base, _ := filepath.Split(dir)
	filename, _ := filepath.Abs(base + "data/sifconfig_test.yaml")
	return ioutil.ReadFile(filename)
}

func LoadTestRaftPersistentStorageFile() ([]byte, error) {
	_, testFile, _, _ := runtime.Caller(0)
	dir, err1 := filepath.Abs(filepath.Dir(testFile))
	if err1 != nil {
		log.Fatal(err1)
	}
	base, _ := filepath.Split(dir)
	filename, _ := filepath.Abs(base + "data/raft_state.json")
	return ioutil.ReadFile(filename)
}

func LoadTestRaftPersistentState() *pb.RaftPersistentState {
	psFile, _ := LoadTestRaftPersistentStorageFile()
	cfg := &pb.RaftPersistentState{}

	err := protojson.Unmarshal(psFile, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func GetCurrentDirPath() string {
	_, testFile, _, _ := runtime.Caller(0)

	base, _ := filepath.Split(filepath.Dir(testFile))
	dir, _ := filepath.Abs(base)
	return dir + "/"
}

func LoadTestRaftConfig() *raftconfig.Config {
	filename, _ := filepath.Abs(GetCurrentDirPath() + "data/sifconfig_test.yaml")

	fileMgr := raftfile.NewFileManager()
	cfg := &raftconfig.Config{}
	file, err2 := fileMgr.LoadFile(filename)

	if err2 != nil {
		log.Fatal(err2)
	}

	err := yaml.Unmarshal(file, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func Get2LogEntries() []*pb.Log {

	logs := []*pb.Log{}
	logs = append(logs, &pb.Log{Term: 1, Message: &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"A": {
				Kind: &structpb.Value_StringValue{
					StringValue: "B",
				},
			}}}},
		&pb.Log{Term: 1, Message: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"C": {
					Kind: &structpb.Value_StringValue{
						StringValue: "B",
					},
				}}}},
	)

	return logs

}
