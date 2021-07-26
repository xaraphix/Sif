package raft

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Raft Node", func() {

	Context("RaftNode initialization", func() {
		When("The Raft node boots up", func() {

			It("Should check if it's booting up from a crash", func() {})
			It("Should load up the logs", func() {})
			It("Should check if it was a leader before crashing", func() {})
			It("Should become a follower if it wasn't a leader on crash", func() {
				node := NewRaftNode(false)
				Expect(node.CurrentRole).To(Equal(FOLLOWER))
			})

			It("Should start the leader heartbeat monitor", func() {
				node := NewRaftNode(false)
				Expect(node.CurrentRole).To(Equal(FOLLOWER))
				time.Sleep(node.LeaderHeartbeatMonitor.TimeoutDuration)
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
			})

		})
	})

	Context("RaftNode Timeouts", func() {
		When("Raft Node's Leader Heartbeat Monitor times out", func() {
			It(`Should become a candidate
			  Vote for iteself
			  Increment the current term
			  request votes from peers`, func() {
				node := NewRaftNode(false)
				term_0 := node.CurrentTerm
				node.LeaderHeartbeatMonitor.LastResetAt = time.Now()

				for {
					if node.CurrentTerm == term_0+1 {
						break
					}
				}

				node.LeaderHeartbeatMonitor.StopMonitor()
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
				Expect(node.VotedFor).To(Equal(raftnode.Id))
				Expect(node.CurrentTerm).To(Equal(term_0 + 1))
				time.Sleep(node.LeaderHeartbeatMonitor.TimeoutDuration * 2)
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
				Expect(node.CurrentTerm).To(Equal(term_0 + 1))

			})
			It("Should Request votes from peers", func() {})
		})

		When("Raft Node's Election times out", func() {
			It(`Should Vote for itself
								 Increment the current term
			           Request votes from peers`, func() {

				node := NewRaftNode(false)
				node.LeaderHeartbeatMonitor.Stopped = true
				term_0 := node.CurrentTerm
				node.StartElectionMonitor()

				//leader heartbeat monitor + election monitor increases term by 2
				// more refined unit test would be stop leader heartbeat monitor when raftNode is created
				for {
					if node.CurrentTerm == term_0+2 {
						break
					}
				}

				node.ElectionMonitor.StopMonitor()
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
				Expect(node.VotedFor).To(Equal(raftnode.Id))
				Expect(node.CurrentTerm).To(Equal(term_0 + 2))
				time.Sleep(node.ElectionMonitor.TimeoutDuration * 2)
				Expect(node.CurrentRole).To(Equal(CANDIDATE))
				Expect(node.CurrentTerm).To(Equal(term_0 + 2))
			})
		})

	})

	Describe("Raft Node Election", func() {
		When("Raft Node Initiates an election", func() {
			It("Should make concurrent RPC calls to all its peers", func() {})
			It("If Majority Votes Against it within the election time duration, it should become a follower", func() {})
			It("If Majority Votes In Favor within the election time duration, it should become a leader", func() {})
			It("On becoming a leader it should send heartbeats to all its peers", func() {})
			It("Should restart election if it cannot make a decision within the election time duration", func() {})
		})
	})
})
