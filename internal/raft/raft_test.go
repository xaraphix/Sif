package raft_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xaraphix/Sif/internal/raft"
	. "github.com/xaraphix/Sif/internal/raft"
	"github.com/xaraphix/Sif/internal/raft/mocks"
	"github.com/xaraphix/Sif/internal/raft/raftelection"
)

var _ = Describe("Raft Node", func() {

	Context("RaftNode initialization", func() {
		When("The Raft node boots up", func() {

			var (
				mockCtrl *gomock.Controller

				node                   *RaftNode
				electionManager        raft.RaftElection
				leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
				electionMonitor        *raft.ElectionMonitor
				rpcAdapter             *mocks.MockRaftRPCAdapter
			)

			BeforeEach(func() {
				mockCtrl = gomock.NewController(GinkgoT())
				electionManager = raftelection.ElctnMgr
				leaderHeartbeatMonitor = LdrHrtbtMontr
				electionMonitor = ElctnMontr
				rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
				rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(VoteResponse{
					VoteGranted: true,
					PeerId:      22,
				}).AnyTimes()

				electionMonitor.Stop()
				leaderHeartbeatMonitor.Stop()
				node = NewRaftNode(electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter)
				node.LeaderHeartbeatMonitor.Start(node)
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
				Expect(node.CurrentRole).To(Equal(FOLLOWER))
			})

			It("Should start the leader heartbeat monitor", func() {
				node.LeaderHeartbeatMonitor.Sleep()
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
			})

		})
	})

	Context("RaftNode Timeouts", func() {
		When("Raft Node's Leader Heartbeat Monitor times out", func() {
			var (
				mockCtrl *gomock.Controller

				node                   *RaftNode
				term_0                 int32
				electionManager        raft.RaftElection
				leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
				electionMonitor        *raft.ElectionMonitor
				rpcAdapter             *mocks.MockRaftRPCAdapter
			)

			BeforeEach(func() {
				mockCtrl = gomock.NewController(GinkgoT())
				electionManager = raftelection.ElctnMgr
				leaderHeartbeatMonitor = LdrHrtbtMontr
				electionMonitor = ElctnMontr
				rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)
				rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(VoteResponse{
					VoteGranted: true,
					PeerId:      int32(22),
				}).AnyTimes()

				electionMonitor.Stop()
				leaderHeartbeatMonitor.Stop()
				node = nil
				node = NewRaftNode(electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter)
				term_0 = node.CurrentTerm
				node.LeaderHeartbeatMonitor.Start(node)
				for {
					if leaderHeartbeatMonitor.GetLastResetAt().IsZero() == false {
						node.LeaderHeartbeatMonitor.Stop()
						break
					}
				}
			})

			It("Should become a candidate", func() {
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
			})

			It("Should Vote for itself", func() {
				Expect(node.VotedFor).To(Equal(node.Id))
			})

			It("Should Increment the current term", func() {
				Expect(node.CurrentTerm).To(BeNumerically("<=", term_0+1))
			})

			It("Should Request votes from peers", func() {
				Expect(node.VotesReceived[0]).To(Equal(int32(22)))

			})
		})

		When("Raft Node's Election times out", func() {
			var (
				mockCtrl *gomock.Controller

				node                   *RaftNode
				term_0                 int32
				electionManager        raft.RaftElection
				leaderHeartbeatMonitor *raft.LeaderHeartbeatMonitor
				electionMonitor        *raft.ElectionMonitor
				rpcAdapter             *mocks.MockRaftRPCAdapter
			)

			BeforeEach(func() {
				mockCtrl = gomock.NewController(GinkgoT())
				electionManager = raftelection.ElctnMgr
				leaderHeartbeatMonitor = LdrHrtbtMontr
				electionMonitor = ElctnMontr
				rpcAdapter = mocks.NewMockRaftRPCAdapter(mockCtrl)

				rpcAdapter.EXPECT().RequestVoteFromPeer(gomock.Any(), gomock.Any()).Return(VoteResponse{
					VoteGranted: true,
					PeerId:      22,
				}).AnyTimes()

				electionMonitor.Stop()
				leaderHeartbeatMonitor.Stop()
				node = nil
				node = NewRaftNode(electionManager, leaderHeartbeatMonitor, electionMonitor, rpcAdapter)
				term_0 = node.CurrentTerm
				node.ElectionMonitor.Start(node)
				for {
					if electionMonitor.GetLastResetAt().IsZero() == false {
						node.ElectionMonitor.Stop()
						break
					}
				}

			})

			It("Should become a candidate", func() {
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
			})

			It("Should Vote for iteself", func() {
				Expect(node.VotedFor).To(Equal(node.Id))
			})

			It("Should Increment the current term", func() {
				//leader heartbeat monitor + election monitor increases term by 2
				// more refined unit test would be stop leader heartbeat monitor when raftNode is created
				Expect(node.CurrentTerm).To(BeNumerically("<=", term_0+1))
			})

			It("Should Request votes from peers", func() {
				Expect(node.VotesReceived[0]).To(Equal(int32(22)))
			})
		})

	})

	Describe("Raft Node Election", func() {
		When("Raft Node Initiates an election", func() {
			It("Should make concurrent RPC calls to all its peers", func() {
				Fail("")
			})
			It("If Majority Votes Against it within the election time duration, it should become a follower", func() {
				Fail("")
			})
			It("If Majority Votes In Favor within the election time duration, it should become a leader", func() {
				Fail("")
			})
			It("On becoming a leader it should send heartbeats to all its peers", func() {
				Fail("")
			})
			It("Should restart election if it cannot make a decision within the election time duration", func() {
				Fail("")
			})
			It("Should become a follower if it discovers a legitimate leader through vote responses", func() {
				Fail("")
			})
			It("Should become a follower if it discovers a legitimate leader through leader heartbeats", func() {
				Fail("")
			})
		})
	})
})
