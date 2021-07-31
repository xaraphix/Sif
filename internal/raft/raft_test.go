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
	"github.com/xaraphix/Sif/internal/raft/raftconfig"
	"github.com/xaraphix/Sif/internal/raft/raftelection"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Sif Raft Consensus", func() {

	Context("RaftNode initialization", func() {

		node := &raft.RaftNode{}
		var setupVars MockSetupVars

		When("The Raft node initializes", func() {
			BeforeEach(func() {
				setupVars := setupRaftNodeInitialization(GinkgoT())
				node = setupVars.node

			})

			AfterEach(func() {
				setupVars.ctrls.configCtrl.Finish()
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

			XIt("should reset currentTerm", func() {
				Fail("Not Yet Implemented")
			})

			XIt("should reset logs", func() {
				Fail("Not Yet Implemented")
			})

			XIt("should reset VotedFor", func() {
				Fail("Not Yet Implemented")
			})

			XIt("should reset CommitLength", func() {
				Fail("Not Yet Implemented")
			})
		})

		When("On booting up from a crash", func() {
			BeforeEach(func() {
				setupVars := setupRaftNodeBootsUpFromCrash(GinkgoT())
				node = setupVars.node
			})

			AfterEach(func() {
				raft.DestructRaftNode(node)
			})

			XIt("should load up currentTerm from persistent storage", func() {
				Fail("Not Yet Implemented")
			})

			XIt("should load up logs from persistent storage", func() {
				Fail("Not Yet Implemented")
			})

			XIt("should load up VotedFor from persistent storage", func() {
				Fail("Not Yet Implemented")
			})

			XIt("should load up CommitLength from persistent storage", func() {
				Fail("Not Yet Implemented")
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
					setupVars = setupLeaderHeartbeatTimeout(GinkgoT())
					node = setupVars.node
					term_0 = setupVars.term_0
				})

				AfterEach(func() {
					setupVars.ctrls.configCtrl.Finish()
					setupVars.ctrls.electionCtrl.Finish()
					setupVars.ctrls.heartCtrl.Finish()
					setupVars.ctrls.rpcCtrl.Finish()
					raft.DestructRaftNode(node)
				})

				It("Should become a candidate", func() {
					Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
				})

				It("Should Vote for itself", func() {
					Expect(node.VotedFor).To(Equal(node.Id))
				})

				It("Should Increment the current term", func() {
					Expect(node.CurrentTerm).To(BeNumerically("==", term_0+1))
				})

				FIt("Should Request votes from peers", func() {
					for {
						if len(node.VotesReceived) == len(node.Peers) {
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
						setupVars.ctrls.configCtrl.Finish()
						setupVars.ctrls.electionCtrl.Finish()
						setupVars.ctrls.heartCtrl.Finish()
						setupVars.ctrls.rpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("Should restart election", func() {
						setupVars := setupRestartElectionOnBeingIndecisive(GinkgoT())
						node = setupVars.node
						term_0 = setupVars.term_0
						loopStartedAt := time.Now()

						for {
							if node.CurrentRole == raft.CANDIDATE &&
								node.CurrentTerm == term_0+1 {
								break
							} else if time.Since(loopStartedAt) > time.Millisecond*600 {
								Fail("Took too much time to be successful")
								break
							}
						}

						loopStartedAt = time.Now()
						for {
							if node.ElectionInProgress == false {
								break
							} else if time.Since(loopStartedAt) > time.Millisecond*900 {
								Fail("Took too much time to be successful")
								break
							}
						}

						loopStartedAt = time.Now()
						for {
							if node.ElectionInProgress == true &&
								node.CurrentTerm == term_0+2 &&
								node.CurrentRole == raft.CANDIDATE {
								break
							} else if time.Since(loopStartedAt) > time.Millisecond*2200 {
								Fail("Took too much time to be successful")
								break
							}
						}

						Succeed()
					})

				})
			})

			Context("Collecting votes", func() {
				var setupVars MockSetupVars
				node := &raft.RaftNode{}
				var sentHeartbeats *map[int]bool
				When("Majority votes in favor", func() {
					AfterEach(func() {
						setupVars.ctrls.configCtrl.Finish()
						setupVars.ctrls.electionCtrl.Finish()
						setupVars.ctrls.heartCtrl.Finish()
						setupVars.ctrls.rpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("should become a leader", func() {
						setupVars := setupMajorityVotesInFavor(GinkgoT())
						node = setupVars.node
						loopStartedAt := time.Now()
						for {
							if node.CurrentRole == raft.LEADER {
								break
							} else if time.Since(loopStartedAt) > time.Millisecond*300 {
								Fail("Took too much time to be successful")
								break
							}

						}
						Succeed()
					})
					XIt("Should cancel the election timer", func() {
						Fail("Not Yet Implemented")
					})

					It("should replicate logs to all its peers", func() {

						config := loadTestRaftConfig()
						setupVars = setupLeaderSendsHeartbeatsOnElectionConclusion(GinkgoT())

						loopStartedAt := time.Now()
						for {
							if setupVars.node.IsHeartBeating == true {
								break
							} else if time.Since(loopStartedAt) > time.Millisecond*300 {
								Fail("Took too much time to be successful")
								break
							}
						}

						node.LeaderHeartbeatMonitor.Sleep()
						Expect((*sentHeartbeats)[int(config.Peers()[0].Id)]).To(Equal(true))
						Expect((*sentHeartbeats)[int(config.Peers()[1].Id)]).To(Equal(true))
					})
				})

				When("The candidate receives a vote response with a higher term than its current term", func() {
					XIt("Should become a follower", func() {

						Fail("Not Yet Implemented")
					})

					XIt("Should cancel election timer", func() {

						Fail("Not Yet Implemented")
					})
				})
			})

			Context("Peers receiving vote requests", func() {
				When("The candidate's log and term are both ok", func() {
					XIt("Should vote yes", func() {

						Fail("Not Yet Implemented")
					})
				})

				When("The candidate's log and term are not ok", func() {
					XIt("Should vote no", func() {

						Fail("Not Yet Implemented")
					})
				})
			})

			Context("Candidate falls back to being a follower", func() {
				node := &raft.RaftNode{}
				When("A legitimate leader sends a heartbeat", func() {
					AfterEach(func() {
						raft.DestructRaftNode(node)
					})
					XIt("Should become a follower", func() {
						Fail("Not Yet Implemented")
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
		When("Leader is replicating log to followers", func() {
			Specify(`The log request should have leaderId, 
			currentTerm, 
			index of the last sent log entry,
			term of the previous log entry,
			commitLength and
			suffix of log entries`, func() {

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

			XIt("Don't know yet what to do", func() {
				Fail("Not Yet Implemented")
			})
		})
	})

})

func setupRaftNode(g GinkgoTInterface, preNodeSetupCB func(
		config raft.RaftConfig,
		file *mocks.MockFile,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
), options SetupOptions) MockSetupVars {

	var (
		node                   *raft.RaftNode
		term_0                 int32
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor

		election   raft.RaftElection
		config     raft.RaftConfig
		file     raftconfig.File
		rpcAdapter raft.RaftRPCAdapter
		heart      raft.RaftHeart
	)

	mockRPCAdapter, rpcCtrl := getMockRPCAdapter(g)
	mockHeart, heartCtrl := getMockHeart(g)
	mockElection, electionCtrl := getMockElection(g)
	mockFile, fileCtrl := getMockFile(g)
	ctrls := Controllers{
		rpcCtrl:      rpcCtrl,
		heartCtrl:    heartCtrl,
		fileCtrl: fileCtrl,
		electionCtrl: electionCtrl,
	}

	election = raftelection.ElctnMgr
	heart = raftelection.LdrHrt

	if options.mockFile {
		file = mockFile
	}

	if options.mockElection {
		election = mockElection
	}

	if options.mockHeart {
		heart = mockHeart
	}

	if options.mockRPCAdapter {
		rpcAdapter = mockRPCAdapter
	}

	config = raftconfig.NewConfig(file)
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	preNodeSetupCB( config, mockFile, mockElection, mockRPCAdapter, mockHeart)

	node = nil
	node = raft.NewRaftNode(config, election, leaderHeartbeatMonitor, rpcAdapter, heart, true)
	term_0 = node.CurrentTerm

	return MockSetupVars{
		node:   node,
		ctrls:  ctrls,
		term_0: term_0,
	}
}

func getMockRPCAdapter(g GinkgoTInterface) (*mocks.MockRaftRPCAdapter, *gomock.Controller) {

	var (
		mockCtrl   *gomock.Controller
		rpcAdapter *mocks.MockRaftRPCAdapter
	)

	mockCtrl = gomock.NewController(g)
	defer mockCtrl.Finish()
	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	return rpcAdapter, mockCtrl
}

func getMockHeart(g GinkgoTInterface) (*mocks.MockRaftHeart, *gomock.Controller) {

	var (
		mockCtrl  *gomock.Controller
		mockHeart *mocks.MockRaftHeart
	)

	mockCtrl = gomock.NewController(g)
	mockHeart = mocks.NewMockRaftHeart(mockCtrl)
	return mockHeart, mockCtrl
}

func getMockConfig(g GinkgoTInterface) (*mocks.MockRaftConfig, *gomock.Controller) {

	var (
		mockCtrl   *gomock.Controller
		mockConfig *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(g)
	mockConfig = mocks.NewMockRaftConfig(mockCtrl)
	return mockConfig, mockCtrl
}

func getMockFile(g GinkgoTInterface) (*mocks.MockFile, *gomock.Controller) {

	var (
		mockCtrl   *gomock.Controller
		mockFile *mocks.MockFile
	)

	mockCtrl = gomock.NewController(g)
	mockFile = mocks.NewMockFile(mockCtrl)
	return mockFile, mockCtrl
}

func getMockElection(g GinkgoTInterface) (*mocks.MockRaftElection, *gomock.Controller) {
	var (
		mockCtrl     *gomock.Controller
		mockElection *mocks.MockRaftElection
	)

	mockCtrl = gomock.NewController(g)
	mockElection = mocks.NewMockRaftElection(mockCtrl)
	return mockElection, mockCtrl
}

type SetupOptions struct {
	mockHeart      bool
	mockElection   bool
	mockFile   bool
	mockConfig     bool
	mockRPCAdapter bool
}

type Controllers struct {
	rpcCtrl      *gomock.Controller
	configCtrl   *gomock.Controller
	fileCtrl   *gomock.Controller
	heartCtrl    *gomock.Controller
	electionCtrl *gomock.Controller
}

type MockSetupVars struct {
	node           *raft.RaftNode
	term_0         int32
	sentHeartbeats *map[int]bool
	ctrls          Controllers
}

func setupRaftNodeBootsUpFromCrash(g GinkgoTInterface) MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockConfig:     true,
		mockRPCAdapter: true,
	}

	preNodeSetupCB := func(
		config raft.RaftConfig,
		file *mocks.MockFile,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
	) {
		file.EXPECT().Load("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		file.EXPECT().Load("./.siflock").AnyTimes().Return(loadTestRaftConfigFile())
		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	return setupRaftNode(g, preNodeSetupCB, options)
}

func setupRaftNodeInitialization(g GinkgoTInterface) MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockConfig:     true,
		mockRPCAdapter: true,
	}

	preNodeSetupCB := func(
		config raft.RaftConfig,
		file *mocks.MockFile,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
	) {
		file.EXPECT().Load("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		file.EXPECT().Load("./.siflock").AnyTimes().Return(nil)
		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	return setupRaftNode(g, preNodeSetupCB, options)
}

func setupLeaderHeartbeatTimeout(g GinkgoTInterface) MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockConfig:     true,
		mockRPCAdapter: true,
	}

	preNodeSetupCB := func(
		config raft.RaftConfig,
		file *mocks.MockFile,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
	) {
		testConfig:= loadTestRaftConfig()
		file.EXPECT().Load("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		file.EXPECT().Load("./.siflock").AnyTimes().Return(nil, errors.New(""))
		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(testConfig.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	setupVars := setupRaftNode(g, preNodeSetupCB, options)
	for {
		if setupVars.node.ElectionInProgress == true {
			setupVars.node.LeaderHeartbeatMonitor.Stop()
			break
		}
	}

	return setupVars
}

func setupMajorityVotesAgainst(g GinkgoTInterface) MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockConfig:     true,
		mockRPCAdapter: true,
	}

	preNodeSetupCB := func(
		config raft.RaftConfig,
		file *mocks.MockFile,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
	) {

		file.EXPECT().Load("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		file.EXPECT().Load("./.siflock").AnyTimes().Return(nil)

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}
	return setupRaftNode(g, preNodeSetupCB, options)
}

func setupMajorityVotesInFavor(g GinkgoTInterface) MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockConfig:     true,
		mockRPCAdapter: true,
	}

	preNodeSetupCB := func(
		config raft.RaftConfig,
		file *mocks.MockFile,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
	) {

		file.EXPECT().Load("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		file.EXPECT().Load("./.siflock").AnyTimes().Return(nil)
		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: true,
			PeerId:      config.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	return setupRaftNode(g, preNodeSetupCB, options)
}

func setupLeaderSendsHeartbeatsOnElectionConclusion(g GinkgoTInterface) MockSetupVars {
	sentHeartbeats := &map[int]bool{
		2: false,
		3: false,
	}

	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockConfig:     true,
		mockRPCAdapter: true,
	}

	preNodeSetupCB := func(
		config raft.RaftConfig,
		file *mocks.MockFile,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
	) {

		file.EXPECT().Load("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		file.EXPECT().Load("./.siflock").AnyTimes().Return(nil)
		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: true,
			PeerId:      config.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[1].Id,
		}).AnyTimes()

		adapter.EXPECT().SendHeartbeatToPeer(config.Peers()[0]).Do(func(interface{}) {
			(*sentHeartbeats)[int(config.Peers()[0].Id)] = true
		}).AnyTimes()

		adapter.EXPECT().SendHeartbeatToPeer(config.Peers()[1]).Do(func(interface{}) {
			(*sentHeartbeats)[int(config.Peers()[1].Id)] = true
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	setupVars := setupRaftNode(g, preNodeSetupCB, options)
	return MockSetupVars{
		node:           setupVars.node,
		term_0:         setupVars.term_0,
		ctrls:          setupVars.ctrls,
		sentHeartbeats: sentHeartbeats,
	}
}

func setupRestartElectionOnBeingIndecisive(g GinkgoTInterface) MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockConfig:     true,
		mockRPCAdapter: true,
	}

	preNodeSetupCB := func(
		config raft.RaftConfig,
		file *mocks.MockFile,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
	) {

		file.EXPECT().Load("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		file.EXPECT().Load("./.siflock").AnyTimes().Return(nil)

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Do(
			func(p raft.Peer, vr raft.VoteRequest) raft.VoteResponse {
				time.Sleep(time.Second)
				return raft.VoteResponse{}
			}).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[1].Id,
		},
		).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	return setupRaftNode(g, preNodeSetupCB, options)
}

func setupGettingLeaderHeartbeatDuringElection(g GinkgoTInterface) MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockConfig:     true,
		mockRPCAdapter: true,
	}

	preNodeSetupCB := func(
		config raft.RaftConfig,
		file *mocks.MockFile,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
	) {

		file.EXPECT().Load("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		file.EXPECT().Load("./.siflock").AnyTimes().Return(nil)
		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}
	return setupRaftNode(g, preNodeSetupCB, options)
}

func setupFindingOtherLeaderThroughVoteResponses(g GinkgoTInterface) MockSetupVars {
	options := SetupOptions{
		mockHeart:      false,
		mockElection:   false,
		mockConfig:     true,
		mockRPCAdapter: true,
	}

	preNodeSetupCB := func(
		config raft.RaftConfig,
		file *mocks.MockFile,
		election *mocks.MockRaftElection,
		adapter *mocks.MockRaftRPCAdapter,
		heart *mocks.MockRaftHeart,
	) {

		file.EXPECT().Load("./sifconfig.yml").AnyTimes().Return(loadTestRaftConfigFile())
		file.EXPECT().Load("./.siflock").AnyTimes().Return(loadTestRaftConfigFile())

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[0].Id,
		}).AnyTimes()

		adapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
			VoteGranted: false,
			PeerId:      config.Peers()[1].Id,
		}).AnyTimes()

		heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	}

	return setupRaftNode(g, preNodeSetupCB, options)
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

func loadTestRaftConfig() raft.RaftConfig {
	_, testFile, _, _ := runtime.Caller(0)
	config := raftconfig.NewConfig(nil)
	cfg := &config
	dir, err1 := filepath.Abs(filepath.Dir(testFile))

	if err1 != nil {
		log.Fatal(err1)
	}
	filename, _ := filepath.Abs(dir + "/mocks/sifconfig_test.yaml")
	yamlFile, _ := ioutil.ReadFile(filename)

	err := yaml.Unmarshal(yamlFile, cfg)

	if err != nil {
		panic(err)
	}
	return config
}
