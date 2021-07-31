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

var _ = Describe("Sif Raft Consensus", func() {

	Context("RaftNode initialization", func() {

		node := &raft.RaftNode{}

		When("The Raft node initializes", func() {

			BeforeEach(func() {
				node = setupRaftNodeBootsUp()
			})

			AfterEach(func() {
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

			XIt("Should check if it's booting up from a crash", func() {
				Fail("Not Yet Implemented")
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
				node = setupRaftNodeBootsUp()
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
			When("Raft Node doesn't receive leader heartbeat for the leader heartbeat duration", func() {

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

		Context("Candidate starts an election", func() {

			Context("Raft Node Election Overtimes", func() {
				When("Candidate is not able to reach to a conclusion within the election allowed time", func() {

					node := &raft.RaftNode{}
					term_0 := int32(0)
					AfterEach(func() {
						raft.DestructRaftNode(node)
					})

					It("Should restart election", func() {
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

				})
			})

			Context("Collecting votes", func() {
				When("Majority votes in favor", func() {
					node := &raft.RaftNode{}
					var sentHeartbeats *map[int]bool
					AfterEach(func() {
						raft.DestructRaftNode(node)
					})

					It("should become a leader", func() {
						node, _ = setupMajorityVotesInFavor()
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
						node, _, sentHeartbeats = setupLeaderSendsHeartbeatsOnElectionConclusion()

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
				XSpecify("Don't know what to do right now", func ()  {
											Fail("Not Yet Implemented")
				})
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
