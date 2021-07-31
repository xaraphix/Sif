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

			It("Should become a follower on booting up", func() {
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
			var sentHeartbeats *map[int]bool
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
				config := loadTestRaftConfig()
				node, term_0, sentHeartbeats = setupLeaderSendsHeartbeatsOnElectionConclusion()

				loopStartedAt := time.Now()
				for {
					if node.IsHeartBeating == true {
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

			It("Should restart election if it cannot make a decision within the election time duration", func() {
				node, term_0 = setupRestartElectionOnBeingIndecisive()

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

			It("Should become a follower if it discovers a legitimate leader through vote responses", func() {
				node, term_0 = setupFindingOtherLeaderThroughVoteResponses()
				Fail("")
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

func setupRaftNode(rpcAdapter *mocks.MockRaftRPCAdapter, heart raft.RaftHeart) (*raft.RaftNode, int32) {

	var (
		mockCtrl *gomock.Controller

		node                   *raft.RaftNode
		term_0                 int32
		electionManager        raft.RaftElection
		leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
		mockConfigFile         *mocks.MockRaftConfig
	)

	mockCtrl = gomock.NewController(GinkgoT())
	electionManager = raftelection.ElctnMgr
	leaderHeartbeatMonitor = raft.NewLeaderHeartbeatMonitor(true)
	mockConfigFile = mocks.NewMockRaftConfig(mockCtrl)

	mockConfigFile.EXPECT().LoadConfig().Return(loadTestRaftConfig()).AnyTimes()

	node = nil
	node = raft.NewRaftNode(mockConfigFile, electionManager, leaderHeartbeatMonitor, rpcAdapter, heart, true)
	term_0 = node.CurrentTerm

	return node, term_0
}

func getMockRPCAdapter() *mocks.MockRaftRPCAdapter {

	var (
		mockCtrl   *gomock.Controller
		rpcAdapter *mocks.MockRaftRPCAdapter
	)

	mockCtrl = gomock.NewController(GinkgoT())
	rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
	return rpcAdapter
}

func getMockHeart() *mocks.MockRaftHeart {

	var (
		mockCtrl  *gomock.Controller
		mockHeart *mocks.MockRaftHeart
	)

	mockCtrl = gomock.NewController(GinkgoT())
	mockHeart = mocks.NewMockRaftHeart(mockCtrl)
	return mockHeart
}

func setupRaftNodeBootsUp() *raft.RaftNode {
	config := loadTestRaftConfig()
	rpcAdapter := getMockRPCAdapter()
	heart := getMockHeart()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	node, _ := setupRaftNode(rpcAdapter, heart)

	return node
}

func setupLeaderHeartbeatTimeout() (*raft.RaftNode, int32) {
	config := loadTestRaftConfig()
	rpcAdapter := getMockRPCAdapter()
	heart := getMockHeart()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	node, term_0 := setupRaftNode(rpcAdapter, heart)
	for {
		if node.ElectionInProgress == true {
			node.LeaderHeartbeatMonitor.Stop()
			break
		}
	}

	return node, term_0
}

func setupMajorityVotesAgainst() (*raft.RaftNode, int32) {
	rpcAdapter := getMockRPCAdapter()
	heart := getMockHeart()
	config := loadTestRaftConfig()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	return setupRaftNode(rpcAdapter, heart)

}

func setupMajorityVotesInFavor() (*raft.RaftNode, int32) {
	rpcAdapter := getMockRPCAdapter()
	heart := getMockHeart()
	config := loadTestRaftConfig()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	return setupRaftNode(rpcAdapter, heart)

}

func setupLeaderSendsHeartbeatsOnElectionConclusion() (*raft.RaftNode, int32, *map[int]bool) {
	sentHeartbeats := &map[int]bool{
		2: false,
		3: false,
	}

	rpcAdapter := getMockRPCAdapter()
	heart := raftelection.LdrHrt
	config := loadTestRaftConfig()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: true,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().SendHeartbeatToPeer(config.Peers()[0]).Do(func(interface{}) {
		(*sentHeartbeats)[int(config.Peers()[0].Id)] = true
	}).AnyTimes()

	rpcAdapter.EXPECT().SendHeartbeatToPeer(config.Peers()[1]).Do(func(interface{}) {
		(*sentHeartbeats)[int(config.Peers()[1].Id)] = true
	}).AnyTimes()

	node, term_0 := setupRaftNode(rpcAdapter, heart)
	return node, term_0, sentHeartbeats
}

func setupRestartElectionOnBeingIndecisive() (*raft.RaftNode, int32) {
	config := loadTestRaftConfig()
	rpcAdapter := getMockRPCAdapter()
	heart := getMockHeart()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Do(
		func(p raft.Peer, vr raft.VoteRequest) raft.VoteResponse {
			time.Sleep(time.Second)
			return raft.VoteResponse{}
		}).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	},
	).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	return setupRaftNode(rpcAdapter, heart)
}

func setupGettingLeaderHeartbeatDuringElection() (*raft.RaftNode, int32) {
	rpcAdapter := getMockRPCAdapter()
	config := loadTestRaftConfig()
	heart := getMockHeart()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	return setupRaftNode(rpcAdapter, heart)

}

func setupFindingOtherLeaderThroughVoteResponses() (*raft.RaftNode, int32) {
	rpcAdapter := getMockRPCAdapter()
	config := loadTestRaftConfig()
	heart := getMockHeart()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[0], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[0].Id,
	}).AnyTimes()

	rpcAdapter.EXPECT().RequestVoteFromPeer(config.Peers()[1], gomock.Any()).Return(raft.VoteResponse{
		VoteGranted: false,
		PeerId:      config.Peers()[1].Id,
	}).AnyTimes()

	heart.EXPECT().StartBeating(gomock.Any()).Return().AnyTimes()
	return setupRaftNode(rpcAdapter, heart)

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
