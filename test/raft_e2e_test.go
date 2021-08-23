package raft_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xaraphix/Sif/internal/raft"
	"github.com/xaraphix/Sif/internal/raft/protos"
	. "github.com/xaraphix/Sif/test/testbed_setup"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ = Describe("Sif Raft Consensus E2E", func() {
	Context("Multiple raft nodes", func() {
		When("A node starts an election", func() {
			Specify("The first node to start the election must win it if there is no other candidate and network errors", func() {

				node1, node2, node3, node4, node5 := Setup5FollowerNodes()

				nodes := make([]**raft.RaftNode, 0)
				nodes = append(nodes, &node1, &node2, &node3, &node4, &node5)
				ProceedWhenRPCAdapterStarted(nodes)

				node1.LeaderHeartbeatMonitor.Start(node1)
				lId := node1.Id

				ProceedWhenLeaderAccepted(nodes, lId)

				Expect(node1.CurrentRole).To(Equal(raft.LEADER))
				Expect(node2.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node3.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node4.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node5.CurrentRole).To(Equal(raft.FOLLOWER))

				cT := node1.CurrentTerm

				Expect(node2.CurrentTerm).To(Equal(cT))
				Expect(node3.CurrentTerm).To(Equal(cT))
				Expect(node4.CurrentTerm).To(Equal(cT))
				Expect(node5.CurrentTerm).To(Equal(cT))

				Expect(node2.CurrentLeader).To(Equal(lId))
				Expect(node3.CurrentLeader).To(Equal(lId))
				Expect(node4.CurrentLeader).To(Equal(lId))
				Expect(node5.CurrentLeader).To(Equal(lId))

				DestructAllNodes(nodes)
				node1, node2, node3, node4, node5 = Setup5FollowerNodes()

				nodes = make([]**raft.RaftNode, 0)
				nodes = append(nodes, &node3, &node1, &node2, &node4, &node5)
				ProceedWhenRPCAdapterStarted(nodes)
				node3.LeaderHeartbeatMonitor.Start(node3)

				lId = node3.Id
				ProceedWhenLeaderAccepted(nodes, lId)

				Expect(node3.CurrentRole).To(Equal(raft.LEADER))
				Expect(node2.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node1.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node4.CurrentRole).To(Equal(raft.FOLLOWER))
				Expect(node5.CurrentRole).To(Equal(raft.FOLLOWER))

				cT = node3.CurrentTerm

				Expect(node2.CurrentTerm).To(Equal(cT))
				Expect(node1.CurrentTerm).To(Equal(cT))
				Expect(node4.CurrentTerm).To(Equal(cT))
				Expect(node5.CurrentTerm).To(Equal(cT))

				Expect(node2.CurrentLeader).To(Equal(lId))
				Expect(node1.CurrentLeader).To(Equal(lId))
				Expect(node4.CurrentLeader).To(Equal(lId))
				Expect(node5.CurrentLeader).To(Equal(lId))

				DestructAllNodes(nodes)
			})

			Specify("If a node becomes a leader and has logs more recent than other nodes, it should replicate the logs to peers", func() {

				node1, node2, node3, node4, node5 := Setup5FollowerNodes()

				node1.Logs = append(node1.Logs, &protos.Log{Term: 1, Message: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"A": {
							Kind: &structpb.Value_StringValue{
								StringValue: "B",
							},
						}}}},
					&protos.Log{Term: 1, Message: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"C": {
								Kind: &structpb.Value_StringValue{
									StringValue: "B",
								},
							}}}},
				)

				nodes := make([]**raft.RaftNode, 0)
				nodes = append(nodes, &node1, &node2, &node3, &node4, &node5)
				ProceedWhenRPCAdapterStarted(nodes)

				Expect(len(node2.Logs)).To(Equal(0))
				Expect(len(node3.Logs)).To(Equal(0))
				Expect(len(node4.Logs)).To(Equal(0))
				Expect(len(node5.Logs)).To(Equal(0))

				node1.LeaderHeartbeatMonitor.Start(node1)
				lId := node1.Id

				ProceedWhenLeaderAccepted(nodes, lId)

				Expect(len(node2.Logs)).To(Equal(2))
				Expect(len(node3.Logs)).To(Equal(2))
				Expect(len(node4.Logs)).To(Equal(2))
				Expect(len(node5.Logs)).To(Equal(2))

			})
		})
	})
})
