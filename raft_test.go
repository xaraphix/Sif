package raft

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Raft Node", func() {

	Context("RaftNode initialization", func() {
		When("The Raft node boots up", func() {

			var node *RaftNode
			BeforeEach(func() {
				node = NewRaftNode(false)
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
				Expect(node.CurrentRole).To(Equal(FOLLOWER))
				Expect(node.ElectionMonitor.Stopped).To(Equal(true))
				time.Sleep(node.LeaderHeartbeatMonitor.TimeoutDuration)
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
			})

		})
	})

	Context("RaftNode Timeouts", func() {
		When("Raft Node's Leader Heartbeat Monitor times out", func() {
			var (
				node   *RaftNode
				term_0 int32
			)

			BeforeEach(func() {
				node = NewRaftNode(false)
				term_0 = node.CurrentTerm

				for {
					if node.CurrentTerm == term_0+1 {
						break
					}
				}
			})

			It("Should become a candidate", func() {
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
			})

			It("Should Vote for itself", func() {
				Expect(node.VotedFor).To(Equal(raftnode.Id))
			})

			It("Should Increment the current term", func() {
				Expect(node.CurrentTerm).To(Equal(term_0 + 1))
			})

			It("Should Request votes from peers", func() {

			})
		})

		When("Raft Node's Election times out", func() {
			var (
				node   *RaftNode
				term_0 int32
			)

			BeforeEach(func() {
				node = NewRaftNode(false)
				term_0 = node.CurrentTerm

				for {
					if node.CurrentTerm == term_0+2 {
						break
					}
				}
				node.ElectionMonitor.StopMonitor()

			})

			It("Should become a candidate", func() {
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
			})

			It("Should Vote for iteself", func() {
				Expect(node.VotedFor).To(Equal(raftnode.Id))
			})

			It("Should Increment the current term", func() {
				//leader heartbeat monitor + election monitor increases term by 2
				// more refined unit test would be stop leader heartbeat monitor when raftNode is created
				Expect(node.CurrentTerm).To(Equal(term_0 + 2))
			})

			It("Should Request votes from peers", func() {
				Fail("")
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
