package raft_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xaraphix/Sif/internal/raft"
	. "github.com/xaraphix/Sif/test/testbed_setup"
)

var _ = Describe("Sif Raft Consensus E2E", func() {
	Context("Multiple raft nodes", func() {
		When("A node starts an election", func() {
			Specify("The first node to start the election must win it if there is no other candidate and network errors", func() {

				node1, node2, node3, node4, node5 := Setup5FollowerNodes()

				node1.LeaderHeartbeatMonitor.Start(node1)

				for e := range node1.GetRaftSignalsChan() {
					if e == raft.BecameLeader {
						break
					}
				}

				Expect(node1.CurrentRole).To(Equal(raft.LEADER))
				Expect(node2.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node3.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node4.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node5.CurrentRole).To(Equal(raft.FOLLOWER))

				raft.DestructRaftNode(node1)
				raft.DestructRaftNode(node2)
				raft.DestructRaftNode(node3)
				raft.DestructRaftNode(node4)
				raft.DestructRaftNode(node5)

				node1, node2, node3, node4, node5 = Setup5FollowerNodes()

				node3.LeaderHeartbeatMonitor.Start(node3)

				for e := range node3.GetRaftSignalsChan() {
					if e == raft.BecameLeader {
						break
					}
				}

				Expect(node3.CurrentRole).To(Equal(raft.LEADER))
				Expect(node2.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node1.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node4.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node5.CurrentRole).To(Equal(raft.FOLLOWER))

				raft.DestructRaftNode(node1)
				raft.DestructRaftNode(node2)
				raft.DestructRaftNode(node3)
				raft.DestructRaftNode(node4)
				raft.DestructRaftNode(node5)

			})
		})
	})
})
