package raft_test

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xaraphix/Sif/internal/raft"
	"github.com/xaraphix/Sif/internal/raft/mocks"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"github.com/xaraphix/Sif/internal/raft/raftadapter"
	"github.com/xaraphix/Sif/internal/raft/raftconfig"
	"github.com/xaraphix/Sif/internal/raft/raftelection"
	"github.com/xaraphix/Sif/internal/raft/raftfile"
	"github.com/xaraphix/Sif/internal/raft/raftlog"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v2"
)

var FileMgr raftfile.RaftFileMgr = raftfile.RaftFileMgr{}
var _ = Describe("Sif Raft Consensus", func() {

	Context("RaftNode initialization", func() {
		node := &raft.RaftNode{}
		var persistentState *pb.RaftPersistentState

		var setupVars MockSetupVars

		When("The Raft node initializes", func() {
			BeforeEach(func() {
				persistentState = loadTestRaftPersistentState()
				setupVars = setupRaftNodeInitialization()
				node = setupVars.node

			})

			AfterEach(func() {
				setupVars.ctrls.fileCtrl.Finish()
				setupVars.ctrls.electionCtrl.Finish()
				setupVars.ctrls.heartCtrl.Finish()
				setupVars.ctrls.rpcCtrl.Finish()
				raft.DestructRaftNode(node)
			})

			It("Should become a follower", func() {
				Expect(node.CurrentRole).To(Equal(raft.FOLLOWER))
			})

			It("Should start the leader heartbeat monitor", func() {
				for {
					if node.ElectionInProgress == true {
						break
					}
				}

				Succeed()
			})

			It("should reset currentTerm", func() {
				Expect(node.CurrentTerm).To(Equal(int32(0)))
			})

			It("should reset logs", func() {
				Expect(len(node.Logs)).To(Equal(0))
			})

			It("should reset VotedFor", func() {
				Expect(node.VotedFor).To(Equal(int32(0)))
			})

			It("should reset CommitLength", func() {
				Expect(node.CommitLength).To(Equal(int32(0)))
			})
		})

		When("On booting up from a crash", func() {
			BeforeEach(func() {
				persistentState = loadTestRaftPersistentState()
				setupVars = setupRaftNodeBootsUpFromCrash()
				node = setupVars.node
			})

			AfterEach(func() {
				setupVars.ctrls.fileCtrl.Finish()
				setupVars.ctrls.electionCtrl.Finish()
				setupVars.ctrls.heartCtrl.Finish()
				setupVars.ctrls.rpcCtrl.Finish()
				raft.DestructRaftNode(node)
			})

			It("should load up currentTerm from persistent storage", func() {
				Expect(node.CurrentTerm).To(Equal(persistentState.RaftCurrentTerm))
			})

			It("should load up logs from persistent storage", func() {
				Expect(len(node.Logs)).To(Equal(len(persistentState.RaftLogs)))
				Expect(len(node.Logs)).To(BeNumerically(">", 0))
			})

			It("should load up VotedFor from persistent storage", func() {
				Expect(node.VotedFor).To(Equal(persistentState.RaftVotedFor))
			})

			It("should load up CommitLength from persistent storage", func() {
				Expect(node.CommitLength).To(Equal(persistentState.RaftCommitLength))
			})
		})
	})

	Context("Raft Node Election", func() {
		Context("RaftNode LeaderHeartbeatMonitor Timeouts", func() {
			node := &raft.RaftNode{}
			term_0 := int32(0)
			var setupVars MockSetupVars

			When("Raft Node doesn't receive leader heartbeat for the leader heartbeat duration", func() {
				BeforeEach(func() {
					setupVars = setupLeaderHeartbeatTimeout()
					node = setupVars.node
					term_0 = setupVars.term_0
				})

				AfterEach(func() {
					defer setupVars.ctrls.fileCtrl.Finish()
					defer setupVars.ctrls.electionCtrl.Finish()
					defer setupVars.ctrls.heartCtrl.Finish()
					defer setupVars.ctrls.rpcCtrl.Finish()
					defer raft.DestructRaftNode(node)
				})

				It("Should become a candidate", func() {
					for e := range node.GetRaftSignalsChan() {
						if e == raft.ElectionTimerStarted {
							break
						}
					}
					Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
				})

				It("Should Vote for itself", func() {
					for e := range node.GetRaftSignalsChan() {
						if e == raft.ElectionTimerStarted {
							break
						}
					}
					Expect(node.VotedFor).To(Equal(node.Id))
				})

				It("Should Increment the current term", func() {
					for e := range node.GetRaftSignalsChan() {
						if e == raft.ElectionTimerStarted {
							break
						}
					}
					Expect(node.CurrentTerm).To(BeNumerically("==", term_0+1))
				})

				It("Should Request votes from peers", func() {
					for {
						if len(node.ElectionMgr.GetReceivedVotes()) == len(node.Peers) {
							break
						}
					}
					Succeed()
				})
			})
		})

		Context("Candidate starts an election", func() {
			Context("Raft Node Election Overtimes", func() {
				node := &raft.RaftNode{}
				term_0 := int32(0)
				var setupVars MockSetupVars
				When("Candidate is not able to reach to a conclusion within the election allowed time", func() {
					AfterEach(func() {
						setupVars.ctrls.fileCtrl.Finish()
						setupVars.ctrls.electionCtrl.Finish()
						setupVars.ctrls.heartCtrl.Finish()
						setupVars.ctrls.rpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("Should restart election", func() {
						setupVars = setupPeerTakesTooMuchTimeToRespond()
						node = setupVars.node
						term_0 = setupVars.term_0

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionRestarted {
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStarted {
								break
							}
						}

						Expect(node.CurrentTerm).To(Equal(term_0 + 2))
					})

				})
			})

			Context("Collecting votes", func() {
				var setupVars MockSetupVars
				node := &raft.RaftNode{}
				When("Majority votes in favor", func() {

					AfterEach(func() {
						setupVars.ctrls.fileCtrl.Finish()
						setupVars.ctrls.electionCtrl.Finish()
						setupVars.ctrls.heartCtrl.Finish()
						setupVars.ctrls.rpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("should become a leader", func() {
						setupVars = setupMajorityVotesInFavor()
						node = setupVars.node
						for e := range node.GetRaftSignalsChan() {
							if e == raft.BecameLeader {
								break
							}
						}
						Succeed()
					})

					It("Should cancel the election timer", func() {
						setupVars = setupLeaderSendsHeartbeatsOnElectionConclusion()
						node = setupVars.node

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStarted {
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStopped {
								break
							}
						}

						counter := 0
						for e := range node.GetRaftSignalsChan() {
							if e == raft.LogRequestSent {
								counter = counter + 1
								if counter == 2 {
									break
								}
							}
						}

						time.Sleep(100 * time.Millisecond)

						Succeed()

					})

					It("should replicate logs to all its peers", func() {

						setupVars = setupLeaderSendsHeartbeatsOnElectionConclusion()
						node = setupVars.node

						for e := range setupVars.node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStarted {
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.BecameLeader {
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.HeartbeatStarted {
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.LogRequestSent {
								Succeed()
								break
							}
						}

					})
				})

				When("The candidate receives a vote response with a higher term than its current term", func() {

					AfterEach(func() {
						setupVars.ctrls.fileCtrl.Finish()
						setupVars.ctrls.electionCtrl.Finish()
						setupVars.ctrls.heartCtrl.Finish()
						setupVars.ctrls.rpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("Should become a follower", func() {
						setupVars = setupCandidateReceivesVoteResponseWithHigherTerm()
						node = setupVars.node
						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStopped {
								break
							}
						}

						Expect(node.CurrentRole).To(Equal(raft.FOLLOWER))
					})

					It("Should Make its current term equal to the one in the voteResponse", func() {
						setupVars = setupCandidateReceivesVoteResponseWithHigherTerm()
						node = setupVars.node
						testConfig := loadTestRaftConfig()
						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStopped {
								break
							}
						}
						Expect(node.CurrentTerm).To(Equal((*setupVars.receivedVoteResponse)[testConfig.Peers()[0].Id].Term))
					})

					It("Should cancel election timer", func() {

						setupVars = setupCandidateReceivesVoteResponseWithHigherTerm()
						node = setupVars.node

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStopped {
								break
							}
						}

						Succeed()
					})
				})
			})

			Context("Peers receiving vote requests", func() {
				var setupVars MockSetupVars
				node := &raft.RaftNode{}

				When("The candidate's log and term are both ok", func() {

					AfterEach(func() {
						setupVars.ctrls.fileCtrl.Finish()
						setupVars.ctrls.electionCtrl.Finish()
						setupVars.ctrls.heartCtrl.Finish()
						setupVars.ctrls.rpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("Should vote yes", func() {
						setupVars = setupPeerReceivingCandidateVoteRequest()
						node = setupVars.node

						node.ElectionMgr.GetResponseForVoteRequest(node, &pb.VoteRequest{
							CurrentTerm: 10,
							NodeId:      0,
							LogLength:   0,
							LastTerm:    0,
						})

						for e := range node.GetRaftSignalsChan() {
							if e == raft.VoteGranted {
								break
							}
						}

						Expect(node.VotedFor).To(Equal(setupVars.leaderId))
					})
				})

				When("The candidate's log and term are not ok", func() {
					AfterEach(func() {
						setupVars.ctrls.fileCtrl.Finish()
						setupVars.ctrls.electionCtrl.Finish()
						setupVars.ctrls.heartCtrl.Finish()
						setupVars.ctrls.rpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("Should vote no", func() {
						setupVars = setupPeerReceivingCandidateVoteRequest()
						node = setupVars.node

						node.ElectionMgr.GetResponseForVoteRequest(node, &pb.VoteRequest{
							CurrentTerm: 0,
							NodeId:      1,
							LogLength:   0,
							LastTerm:    10,
						})

						for e := range node.GetRaftSignalsChan() {
							if e == raft.VoteNotGranted {
								break
							}
						}

						Expect(node.VotedFor).To(Equal(int32(0)))
						raft.DestructRaftNode(node)
					})
				})
			})

			Context("Candidate falls back to being a follower", func() {
				var setupVars MockSetupVars
				node := &raft.RaftNode{}
				When("A legitimate leader sends a heartbeat", func() {
					AfterEach(func() {
						setupVars.ctrls.fileCtrl.Finish()
						setupVars.ctrls.electionCtrl.Finish()
						setupVars.ctrls.heartCtrl.Finish()
						setupVars.ctrls.rpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})
					It("Should become a follower", func() {
						setupVars = setupMajorityVotesAgainst()
						node = setupVars.node

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStarted {
								node.ElectionMgr.GetLeaderHeartChannel() <- raft.RaftNode{
									Node: raft.Node{
										CurrentTerm: int32(10),
									},
								}
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.BecameFollower {
								break
							}
						}

						Expect(node.CurrentRole).To(Equal(raft.FOLLOWER))
						Expect(node.CurrentTerm).To(Equal(int32(10)))
					})
				})
			})

			When("Candidate receives vote request from another candidate", func() {
				XSpecify("Don't know what to do right now", func() {
					Fail("Not Yet Implemented")
				})
			})
		})
	})

	Context("Broadcasting Messages", func() {
		When("A Leader receives a broadcast request", func() {

			XIt("Should Append the entry to its log", func() {

			})

			XIt("Should update the Acknowledged log length for itself", func() {

			})

			XIt("Should Replicate the log to its followers after appending to its own logs", func() {

			})
		})

		When("A Raft node is a leader", func() {
			XIt("Should send heart beats to its followers", func() {

			})
		})

		When("A Follower receives a broadcast request", func() {
			XIt("Should forward the request to the leader node", func() {

			})
		})
	})

	Context("Log Replication", func() {
		var setupVars MockSetupVars
		node := &raft.RaftNode{}
		When("Leader is replicating log to followers", func() {

			AfterEach(func() {
				setupVars.ctrls.fileCtrl.Finish()
				setupVars.ctrls.electionCtrl.Finish()
				setupVars.ctrls.heartCtrl.Finish()
				setupVars.ctrls.rpcCtrl.Finish()
				raft.DestructRaftNode(node)
			})

			Specify(`The log request should have leaderId, 
			currentTerm, 
			index of the last sent log entry,
			term of the previous log entry,
			commitLength and
			suffix of log entries`, func() {

				setupVars = setupMajorityVotesInFavor()
				node = setupVars.node

				for e := range node.GetRaftSignalsChan() {
					if e == raft.LogRequestSent {
						break
					}
				}

				Expect((*setupVars.sentLogReplicationReq).CurrentTerm).To(Equal(node.CurrentTerm))
				Succeed()
			})
		})

		When("Follower receives a log replication request", func() {
			When("the received log is ok", func() {

				XIt("Should update its currentTerm to to new term", func() {

				})

				XIt("Should continue being the follower and cancel the election it started (if any)", func() {

				})

				XIt("Should update its acknowledged length of the log", func() {

				})

				XIt("Should append the entry to its log", func() {

				})

				XIt("Should send the log response back to the leader with the acknowledged new length", func() {

				})

			})

			When("The received log or term is not ok", func() {

				XIt(`Should send its currentTerm
			acknowledged length as 0
			acknowledgment as false to the leader`, func() {

				})
			})

		})
	})

	Context("A node updates its log based on the log received from leader", func() {
		//TODO
	})

	Context("Leader Receives log acknowledgments", func() {
		When("The term in the log acknowledgment is more than leader's currentTerm", func() {})
		When("The term in the log acknowledgment is ok and log replication has been acknowledged by the follower", func() {

			When("log replication is acknowledged by the follower", func() {

				XIt(`Should update its acknowledged and 
					sentLength for the follower`, func() {

				})

				XIt("Should commit the log entries to persistent storage", func() {

				})
			})

			When("log replication is not acknowledged by the follower", func() {

				XIt("Should update the sent length for the follower to one less than the previously attempted sent length value of the log", func() {

				})
				XIt("Should send a log replication request to the follower with log length = last request log length - 1", func() {

				})
			})

		})

	})

	Context("Commiting log entries to persistent storage", func() {
		//TODO
		When("The commit is successful it should send the log message to the client of Sif", func() {

			FIt("Don't know yet what to do", func() {

				node1, node2, node3, node4, node5 := setup5FollowerNodes()

				if node1 == nil || node2 == nil || node3 == nil || node4 == nil || node5 == nil {

				}

				for {
					time.Sleep(5 * time.Second)
				}
			})
		})
	})

})

func setupRaftNode(preNodeSetupCB func(
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

	mockRPCAdapter, rpcCtrl := getMockRPCAdapter()
	mockHeart, heartCtrl := getMockHeart()
	mockElection, electionCtrl := getMockElection()
	mockFile, fileCtrl := getMockFile()
	mockLog, logCtrl := getMockLog()
	mockMonitor, monitorCtrl := getMockMonitor()

	ctrls := Controllers{
		rpcCtrl:      rpcCtrl,
		logCtrl:      logCtrl,
		heartCtrl:    heartCtrl,
		fileCtrl:     fileCtrl,
		electionCtrl: electionCtrl,
		monitorCtrl:  monitorCtrl,
	}

	raftOptions := raft.RaftOptions{}

	if options.startMonitor {
		raftOptions.StartLeaderHeartbeatMonitorAfterInitializing = true
	}

	election = raftelection.NewElectionManager()
	heart = raftelection.NewLeaderHeart()
	logMgr = &raftlog.LogMgr{}
	monitor = raft.NewLeaderHeartbeatMonitor(true)

	if options.mockFile {
		fileMgr = mockFile
	}

	if options.mockElection {
		election = mockElection
	}

	if options.mockHeart {
		heart = mockHeart
	}

	if options.mockLeaderMonitor {
		monitor = mockMonitor
	}

	if options.mockRPCAdapter {
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
		node:   node,
		ctrls:  ctrls,
		term_0: term_0,
	}
}

func getMockRPCAdapter() (*mocks.MockRaftRPCAdapter, *gomock.Controller) {

	var (
		mockCtrl   *gomock.Controller
		rpcAdapter *mocks.MockRaftRPCAdapter
	)

	mockCtrl = gomock.NewController(GinkgoT())
	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	return rpcAdapter, mockCtrl
}

func getMockHeart() (*mocks.MockRaftHeart, *gomock.Controller) {

	var (
		mockCtrl  *gomock.Controller
		mockHeart *mocks.MockRaftHeart
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockHeart = mocks.NewMockRaftHeart(mockCtrl)
	return mockHeart, mockCtrl
}

func getMockConfig() (*mocks.MockRaftConfig, *gomock.Controller) {

	var (
		mockCtrl   *gomock.Controller
		mockConfig *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockConfig = mocks.NewMockRaftConfig(mockCtrl)
	return mockConfig, mockCtrl
}

func getMockFile() (*mocks.MockRaftFile, *gomock.Controller) {

	var (
		mockCtrl *gomock.Controller
		mockFile *mocks.MockRaftFile
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockFile = mocks.NewMockRaftFile(mockCtrl)
	return mockFile, mockCtrl
}

func getMockLog() (*mocks.MockRaftLog, *gomock.Controller) {

	var (
		mockCtrl *gomock.Controller
		mockLog  *mocks.MockRaftLog
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockLog = mocks.NewMockRaftLog(mockCtrl)
	return mockLog, mockCtrl
}

func getMockMonitor() (*mocks.MockRaftMonitor, *gomock.Controller) {

	var (
		mockCtrl *gomock.Controller
		mockLog  *mocks.MockRaftMonitor
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockLog = mocks.NewMockRaftMonitor(mockCtrl)
	return mockLog, mockCtrl
}

func getMockElection() (*mocks.MockRaftElection, *gomock.Controller) {
	var (
		mockCtrl     *gomock.Controller
		mockElection *mocks.MockRaftElection
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockElection = mocks.NewMockRaftElection(mockCtrl)
	return mockElection, mockCtrl
}

type SetupOptions struct {
	mockHeart         bool
	mockLog           bool
	mockElection      bool
	mockFile          bool
	mockConfig        bool
	mockRPCAdapter    bool
	mockLeaderMonitor bool
	startMonitor      bool
}

type Controllers struct {
	rpcCtrl      *gomock.Controller
	logCtrl      *gomock.Controller
	configCtrl   *gomock.Controller
	fileCtrl     *gomock.Controller
	heartCtrl    *gomock.Controller
	electionCtrl *gomock.Controller
	monitorCtrl  *gomock.Controller
}

type MockSetupVars struct {
	node                  *raft.RaftNode
	term_0                int32
	sentHeartbeats        *map[int]bool
	sentVoteRequests      *map[int]*pb.VoteRequest
	receivedVoteResponse  *map[int32]*pb.VoteResponse
	sentLogReplicationReq **pb.LogRequest
	ctrls                 Controllers
	leaderId              int32
}

func setupRaftNodeBootsUpFromCrash() MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockLog:        true,
		mockFile:       true,
		mockRPCAdapter: true,
		startMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {
		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, nil)
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))

		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
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

	return setupRaftNode(preNodeSetupCB, options)
}

func setupRaftNodeInitialization() MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockFile:       true,
		mockLog:        true,
		mockRPCAdapter: true,
		startMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {
		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()

		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
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

	return setupRaftNode(preNodeSetupCB, options)
}

func setupLeaderHeartbeatTimeout() MockSetupVars {
	options := SetupOptions{
		mockFile:       true,
		mockHeart:      false,
		mockLog:        true,
		mockElection:   false,
		mockRPCAdapter: true,
		startMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {
		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
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

	setupVars := setupRaftNode(preNodeSetupCB, options)
	for {
		if setupVars.node.ElectionInProgress == true {
			setupVars.node.LeaderHeartbeatMonitor.Stop()
			break
		}
	}

	return setupVars
}

func setupMajorityVotesAgainst() MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockFile:       true,
		mockLog:        true,
		mockRPCAdapter: true,
		startMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()

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
	return setupRaftNode(preNodeSetupCB, options)
}

func setupMajorityVotesInFavor() MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockFile:       true,
		mockLog:        true,
		mockRPCAdapter: true,
		startMonitor:   true,
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

		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
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

	setupVars := setupRaftNode(preNodeSetupCB, options)
	setupVars.sentLogReplicationReq = &logReplicationReq
	return setupVars
}

func setupLeaderSendsHeartbeatsOnElectionConclusion() MockSetupVars {
	sentHeartbeats := &map[int]bool{
		2: false,
		3: false,
	}

	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockFile:       true,
		mockLog:        true,
		mockRPCAdapter: true,
		startMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
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

	setupVars := setupRaftNode(preNodeSetupCB, options)
	return MockSetupVars{
		node:           setupVars.node,
		term_0:         setupVars.term_0,
		ctrls:          setupVars.ctrls,
		sentHeartbeats: sentHeartbeats,
	}
}

func setupPeerTakesTooMuchTimeToRespond() MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockFile:       true,
		mockLog:        true,
		mockRPCAdapter: true,
		startMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()

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

	return setupRaftNode(preNodeSetupCB, options)
}

func setupGettingLeaderHeartbeatDuringElection() MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockLog:        true,
		mockFile:       true,
		mockRPCAdapter: true,
		startMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
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
	return setupRaftNode(preNodeSetupCB, options)
}

func setupFindingOtherLeaderThroughVoteResponses() MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockFile:       true,
		mockLog:        true,
		mockRPCAdapter: true,
		startMonitor:   true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {
		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))

		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
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

	return setupRaftNode(preNodeSetupCB, options)
}

func setupCandidateReceivesVoteResponseWithHigherTerm() MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockLog:        true,
		mockFile:       true,
		mockRPCAdapter: true,
		startMonitor:   true,
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

		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()
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
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Return(vr1).AnyTimes()
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(vr2).AnyTimes()
		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()

		sentVoteResponse[testConfig.Peers()[0].Id] = vr1
		sentVoteResponse[testConfig.Peers()[1].Id] = vr2
	}

	setupVars := setupRaftNode(preNodeSetupCB, options)
	setupVars.receivedVoteResponse = &sentVoteResponse
	return setupVars
}

func setupPeerReceivingCandidateVoteRequest() MockSetupVars {

	options := SetupOptions{
		mockHeart:         false,
		mockElection:      false,
		mockLog:           false,
		mockFile:          true,
		mockRPCAdapter:    true,
		mockLeaderMonitor: true,
	}

	preNodeSetupCB := func(
		fileMgr *mocks.MockRaftFile,
		logMgr *mocks.MockRaftLog,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
		monitor *mocks.MockRaftMonitor,
	) {

		testConfig := loadTestRaftConfig()
		testPersistentStorageFile, _ := loadTestRaftPersistentStorageFile()
		adapter.EXPECT().StartAdapter(gomock.Any()).Return().AnyTimes()
		fileMgr.EXPECT().LoadFile("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+".siflock").AnyTimes().Return(nil, errors.New(""))
		fileMgr.EXPECT().LoadFile(testConfig.RaftInstanceDirPath+"raft_state.json").AnyTimes().Return(testPersistentStorageFile, errors.New(""))
	}

	setupVars := setupRaftNode(preNodeSetupCB, options)

	return setupVars
}

func getConfig() raftconfig.Config {
	return raftconfig.Config{}
}

func loadTestRaftConfigFile() ([]byte, error) {
	_, testFile, _, _ := runtime.Caller(0)
	dir, err1 := filepath.Abs(filepath.Dir(testFile))
	if err1 != nil {
		log.Fatal(err1)
	}
	filename, _ := filepath.Abs(dir + "/mocks/sifconfig_test.yaml")
	return ioutil.ReadFile(filename)
}

func loadTestRaftPersistentStorageFile() ([]byte, error) {
	_, testFile, _, _ := runtime.Caller(0)
	dir, err1 := filepath.Abs(filepath.Dir(testFile))
	if err1 != nil {
		log.Fatal(err1)
	}
	filename, _ := filepath.Abs(dir + "/mocks/raft_state.json")
	return ioutil.ReadFile(filename)
}

func loadTestRaftPersistentState() *pb.RaftPersistentState {
	psFile, _ := loadTestRaftPersistentStorageFile()
	cfg := &pb.RaftPersistentState{}

	err := protojson.Unmarshal(psFile, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func getCurrentDirPath() string {
	_, testFile, _, _ := runtime.Caller(0)

	dir, _ := filepath.Abs(filepath.Dir(testFile))
	return dir
}

func loadTestRaftConfig() *raftconfig.Config {
	filename, _ := filepath.Abs(getCurrentDirPath() + "/mocks/sifconfig_test.yaml")

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

func setup5FollowerNodes() (*raft.RaftNode, *raft.RaftNode, *raft.RaftNode, *raft.RaftNode, *raft.RaftNode) {

	currDirPath := getCurrentDirPath()

	deps1 := getNewNodeDeps()
	deps1.ConfigManager.SetConfigFilePath(currDirPath + "/mocks/node1_config.yml")
	node1 := raft.NewRaftNode(deps1)

	deps2 := getNewNodeDeps()
	deps2.ConfigManager.SetConfigFilePath(currDirPath + "/mocks/node2_config.yml")
	node2 := raft.NewRaftNode(deps2)

	deps3 := getNewNodeDeps()
	deps3.ConfigManager.SetConfigFilePath(currDirPath + "/mocks/node3_config.yml")
	node3 := raft.NewRaftNode(deps3)

	deps4 := getNewNodeDeps()
	deps4.ConfigManager.SetConfigFilePath(currDirPath + "/mocks/node4_config.yml")
	node4 := raft.NewRaftNode(deps4)

	deps5 := getNewNodeDeps()
	deps5.ConfigManager.SetConfigFilePath(currDirPath + "/mocks/node5_config.yml")
	node5 := raft.NewRaftNode(deps5)

	return node1, node2, node3, node4, node5
}

func getNewNodeDeps() raft.RaftDeps {
	options := raft.RaftOptions{
		StartLeaderHeartbeatMonitorAfterInitializing: true,
	}

	return raft.RaftDeps{
		FileManager:      raftfile.NewFileManager(),
		ConfigManager:    raftconfig.NewConfig(),
		ElectionManager:  raftelection.NewElectionManager(),
		HeartbeatMonitor: raft.NewLeaderHeartbeatMonitor(true),
		RPCAdapter:       raftadapter.NewRaftNodeAdapter(),
		LogManager:       raftlog.NewLogManager(),
		Heart:            raftelection.NewLeaderHeart(),
		Options:          options,
	}

}
