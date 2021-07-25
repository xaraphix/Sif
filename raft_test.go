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
				time.Sleep(node.LeaderHeartbeatMonitor.HeartbeatTimeout)
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
				node.LeaderHeartbeatMonitor.LeaderLastHeartbeat = time.Now()
				term_0 := node.CurrentTerm
				time.Sleep(node.LeaderHeartbeatMonitor.HeartbeatTimeout)

				Expect(node.CurrentRole).To(Equal(CANDIDATE))
				Expect(node.VotedFor).To(Equal(raftnode.Id))
				Expect(node.CurrentTerm).To(Equal(term_0 + 1))
			})
			It("Should Request votes from peers", func() {})
		})

		When("Raft Node's Election times out", func() {
			It("Should Vote for itself", func() {})
			It("Should Increment the current term", func() {})
			It("Should Request votes from peers", func() {})
		})

	})
})
