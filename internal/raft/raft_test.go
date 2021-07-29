package raft_test

import (
	"io/ioutil"
	"path/filepath"
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
				node.LeaderHeartbeatMonitor.Sleep()
				Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
			})

		})
	})

	Context("RaftNode Timeouts", func() {
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
				Expect(node.VotesReceived[0]).To(Equal(int32(22)))

			})
		})

		When("Raft Node's Election times out", func() {

			node := &raft.RaftNode{}
			term_0 := int32(0)

			BeforeEach(func() {
				node, term_0 = setupElectionTimerTimout()
			})

			AfterEach(func() {
				raft.DestructRaftNode(node)
			})

			It("Should become a candidate", func() {
				Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
			})

			It("Should Vote for iteself", func() {
				Expect(node.VotedFor).To(Equal(node.Id))
			})

			It("Should Increment the current term", func() {
				//leader heartbeat monitor + election monitor increases term by 2
				// more refined unit test would be stop leader heartbeat monitor when raftNode is created
				Expect(node.CurrentTerm).To(BeNumerically("==", term_0+1))
			})

			It("Should Request votes from peers", func() {
				for {
					if len(node.VotesReceived) == len(node.Peers) {
						break
					}
				}
				Expect(node.VotesReceived[0]).To(Equal(int32(22)))
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

			It("If Majority Votes Against it within the election time duration, it should become a follower", func() {
				node, term_0 = setupMajorityVotesAgainst()
				for {
					if len(node.VotesReceived) == len(node.Peers) {
						break
					}
				}
				Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
				Expect(node.CurrentTerm).To(Equal(term_0 + 1))
				time.Sleep(100 * time.Millisecond)
				Expect(node.CurrentRole).To(Equal(raft.FOLLOWER))

			})

			It("If Majority Votes In Favor within the election time duration, it should become a leader", func() {
				node, term_0 = setupMajorityVotesInFavor()
				Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
				time.Sleep(100 * time.Millisecond)
				Expect(node.CurrentRole).To(Equal(raft.LEADER))
			})

			It("On becoming a leader it should send heartbeats to all its peers", func() {
				node, term_0 = setupMajorityVotesInFavor()
				Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
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
		electionMonitor        *raft.ElectionMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	electionMonitor = raft.NewElectionMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      22,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	electionMonitor.Stop()
	leaderHeartbeatMonitor.Stop()
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter, true)
	node.LeaderHeartbeatMonitor.Start(node)

	return node
}

func setupLeaderHeartbeatTimeout() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		electionMonitor        *raft.ElectionMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	electionMonitor = raft.NewElectionMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      22,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	electionMonitor.Stop()
	leaderHeartbeatMonitor.Stop()
	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm
	node.LeaderHeartbeatMonitor.Start(node)
	for {
		if node.ElectionInProgress == true {
			node.LeaderHeartbeatMonitor.Stop()
			break
		}
	}

	return node, term_0
}

func setupElectionTimerTimout() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		electionMonitor        *raft.ElectionMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	electionMonitor = raft.NewElectionMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      22,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	electionMonitor.Stop()
	leaderHeartbeatMonitor.Stop()
	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm
	node.ElectionMonitor.Start(node)
	for {
		if node.ElectionInProgress == true {
			node.ElectionMonitor.Stop()
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
		electionMonitor        *raft.ElectionMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	electionMonitor = raft.NewElectionMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      22,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	electionMonitor.Stop()
	leaderHeartbeatMonitor.Stop()
	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm
	node.ElectionMonitor.Start(node)
	for {
		if node.ElectionInProgress == true {
			node.ElectionMonitor.Stop()
			break
		}
	}

	return node, term_0

}
func setupMajorityVotesInFavor() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		electionMonitor        *raft.ElectionMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	electionMonitor = raft.NewElectionMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      22,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	electionMonitor.Stop()
	leaderHeartbeatMonitor.Stop()
	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm
	node.ElectionMonitor.Start(node)
	for {
		if node.ElectionInProgress == true {
			node.ElectionMonitor.Stop()
			break
		}
	}

	return node, term_0

}
func setupRestartElectionOnBeingIndecisive() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		electionMonitor        *raft.ElectionMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	electionMonitor = raft.NewElectionMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      22,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	electionMonitor.Stop()
	leaderHeartbeatMonitor.Stop()
	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm
	node.ElectionMonitor.Start(node)
	for {
		if node.ElectionInProgress == true {
			node.ElectionMonitor.Stop()
			break
		}
	}

	return node, term_0

}
func setupGettingLeaderHeartbeatDuringElection() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		electionMonitor        *raft.ElectionMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	electionMonitor = raft.NewElectionMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      22,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	electionMonitor.Stop()
	leaderHeartbeatMonitor.Stop()
	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm
	node.ElectionMonitor.Start(node)
	for {
		if node.ElectionInProgress == true {
			node.ElectionMonitor.Stop()
			break
		}
	}

	return node, term_0

}
func setupFindingOtherLeaderThroughVoteResponses() (*raft.RaftNode, int32) {
	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		electionMonitor        *raft.ElectionMonitor
		rpcAdapter             *mocks.MockRaftRPCAdapter
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	electionMonitor = raft.NewElectionMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      22,
	}).AnyTimes()

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	electionMonitor.Stop()
	leaderHeartbeatMonitor.Stop()
	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter, true)
	term_0 = node.CurrentTerm
	node.ElectionMonitor.Start(node)
	for {
		if node.ElectionInProgress == true {
			node.ElectionMonitor.Stop()
			break
		}
	}

	return node, term_0

}

func getConfig() raftconfig.Config {
	return raftconfig.Config{}
}

func loadTestRaftConfig() raft.RaftConfig {

	config := raftconfig.NewConfig()
	cfg := &config
	filename, _ := filepath.Abs("./mocks/sifconfig_test.yaml")
	yamlFile, _ := ioutil.ReadFile(filename)

	err := yaml.Unmarshal(yamlFile, cfg)

	if err != nil {
		panic(err)
	}
	return config
}
