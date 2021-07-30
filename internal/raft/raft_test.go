package raft_test

import (
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

var _ = Describe("Raft Node", func() {

	Context("RaftNode initialization", func() {
		When("The Raft node boots up", func() {
			node := &raft.RaftNode{}
			BeforeEach(func() {
				node = setupRaftNodeBootsUp()
			})

			AfterEach(func() {
				raft.DestructRaftNode(node)
			})

			It("Should check if it's booting up from a crash", func() {
				Fail("")
			})

			It("Should load up the logs", func() {
				Fail("")
			})

			It("Should check if it was a leader before crashing", func() {
				Fail("")
			})

			It("Should become a follower if it wasn't a leader on crash", func() {
				Fail("")
			})

			It("Should become a follower if it is not booting up from a crash", func() {
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

		})
	})

	Context("RaftNode LeaderHeartbeatMonitor Timeouts", func() {
		When("Raft Node's Leader Heartbeat Monitor times out", func() {
			node := &raft.RaftNode{}
			term_0 := int32(0)

			BeforeEach(func() {
				node, term_0 = setupLeaderHeartbeatTimeout()
			})

			AfterEach(func() {
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

			It("Should Request votes from peers", func() {
				for {
					if len(node.VotesReceived) == len(node.Peers) {
						break
					}
				}

				Succeed()
			})
		})

	})

	Describe("Raft Node Election", func() {
		When("Raft Node Initiates an election", func() {

			node := &raft.RaftNode{}
			term_0 := int32(0)

			AfterEach(func() {
				raft.DestructRaftNode(node)
			})

			It("should become a follower If Majority Votes Against it within the election time duration", func() {
				node, term_0 = setupMajorityVotesAgainst()
				loopStartedAt := time.Now()
				for {
					if node.CurrentRole == raft.FOLLOWER && node.CurrentTerm >= term_0+1 {
						break
					} else if time.Since(loopStartedAt) > time.Millisecond*300 {
						Fail("Took too much time to be successful")
						break
					}
				}

				Succeed()
			})

			It("should become a leader If Majority Votes In Favor within the election time duration ", func() {
				node, term_0 = setupMajorityVotesInFavor()
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

			It("should send heartbeats to all its peers On becoming a leader ", func() {
				node, term_0 = setupMajorityVotesInFavor()
				Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
				Expect(node.CurrentTerm).To(Equal(term_0 + 1))
				for {
					if len(node.VotesReceived) == len(node.Peers) {
						break
					}
				}
				time.Sleep(100 * time.Millisecond)
				Expect(node.CurrentRole).To(Equal(raft.LEADER))

			})
			It("Should restart election if it cannot make a decision within the election time duration", func() {
				node, term_0 = setupRestartElectionOnBeingIndecisive()
				Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
				time.Sleep(100 * time.Millisecond)
				Expect(node.CurrentRole).To(Equal(raft.LEADER))
			})
			It("Should become a follower if it discovers a legitimate leader through vote responses", func() {
				node, term_0 = setupFindingOtherLeaderThroughVoteResponses()
				Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
				time.Sleep(100 * time.Millisecond)
				Expect(node.CurrentRole).To(Equal(raft.LEADER))
			})
			It("Should become a follower if it discovers a legitimate leader through leader heartbeats", func() {
				node, term_0 = setupGettingLeaderHeartbeatDuringElection()
				Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
				time.Sleep(100 * time.Millisecond)
				Expect(node.CurrentRole).To(Equal(raft.LEADER))
			})
		})
	})
})

func setupRaftNodeBootsUp() *raft.RaftNode {

	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	config := loadTestRaftConfig()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, rpcAdapter, true)

	return node
}

func setupLeaderHeartbeatTimeout() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	config := loadTestRaftConfig()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm
	for {
		if node.ElectionInProgress == true {
			node.LeaderHeartbeatMonitor.Stop()
			break
		}
	}

	return node, term_0
}
func setupMajorityVotesAgainst() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)

	config := loadTestRaftConfig()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm

	return node, term_0

}
func setupMajorityVotesInFavor() (*raft.RaftNode, int32) {

	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)

	config := loadTestRaftConfig()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm

	return node, term_0

}
func setupRestartElectionOnBeingIndecisive() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	config := loadTestRaftConfig()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).DoAndReturn(
		func(p raft.Peer, vr raft.VoteRequest) raft.VoteResponse {
			// checks whatever

			return raft.VoteResponse{
				VoteGranted: false,
				PeerId:      config.Peers()[1].Id,
			}
		}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm

	return node, term_0
}
func setupGettingLeaderHeartbeatDuringElection() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	config := loadTestRaftConfig()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm
	return node, term_0

}
func setupFindingOtherLeaderThroughVoteResponses() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)

	config := loadTestRaftConfig()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm

	return node, term_0

}

func getConfig() raftconfig.Config {
	return raftconfig.Config{}
}

func loadTestRaftConfig() raft.RaftConfig {
	_, testFile, _, _ := runtime.Caller(0)
	config := raftconfig.NewConfig()
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
