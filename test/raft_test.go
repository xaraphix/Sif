package raft_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xaraphix/Sif/internal/raft"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"github.com/xaraphix/Sif/internal/raft/raftfile"
	. "github.com/xaraphix/Sif/test/testbed_setup"
)

var FileMgr raftfile.RaftFileMgr = raftfile.RaftFileMgr{}
var _ = Describe("Sif Raft Consensus", func() {
	Context("RaftNode initialization", func() {
		node := &raft.RaftNode{}
		var persistentState *pb.RaftPersistentState

		var setupVars MockSetupVars

		When("The Raft node initializes", func() {
			BeforeEach(func() {
				persistentState = LoadTestRaftPersistentState()
				setupVars = SetupRaftNodeInitialization()
				node = setupVars.Node

			})

			AfterEach(func() {
				setupVars.Ctrls.FileCtrl.Finish()
				setupVars.Ctrls.ElectionCtrl.Finish()
				setupVars.Ctrls.HeartCtrl.Finish()
				setupVars.Ctrls.RpcCtrl.Finish()
				defer node.Close()
			})

			It("Should become a follower", func() {
				Expect(node.CurrentRole).To(Equal(raft.FOLLOWER))
			})

			It("Should start the leader heartbeat monitor", func() {
				CheckIfEventTriggered(node, raft.BecameCandidate, raft.RaftEventDetails{})
				Succeed()
			})

			It("should reset currentTerm", func() {
				Expect(node.CurrentTerm).To(Equal(int32(0)))
			})

			It("should reset logs", func() {
				Expect(len(node.Logs)).To(Equal(0))
			})

			It("should reset VotedFor", func() {
				Expect(node.VotedFor).To(Equal(""))
			})

			It("should reset CommitLength", func() {
				Expect(node.CommitLength).To(Equal(int32(0)))
			})
		})

		When("On booting up from a crash", func() {
			BeforeEach(func() {
				persistentState = LoadTestRaftPersistentState()
				setupVars = SetupRaftNodeBootsUpFromCrash()
				node = setupVars.Node
			})

			AfterEach(func() {
				setupVars.Ctrls.FileCtrl.Finish()
				setupVars.Ctrls.ElectionCtrl.Finish()
				setupVars.Ctrls.HeartCtrl.Finish()
				setupVars.Ctrls.RpcCtrl.Finish()
				defer node.Close()
			})

			It("should load up currentTerm from persistent storage", func() {
				Expect(node.CurrentTerm).To(Equal(persistentState.RaftCurrentTerm))
			})

			It("should load up logs from persistent storage", func() {
				Expect(len(node.Logs)).To(Equal(len(persistentState.RaftLogs)))
				Expect(len(node.Logs)).To(BeNumerically(">", 0))
			})

			It("should load up VotedFor from persistent storage", func() {
				Expect(node.VotedFor).To(Equal(persistentState.RaftVotedFor))
			})

			It("should load up CommitLength from persistent storage", func() {
				Expect(node.CommitLength).To(Equal(persistentState.RaftCommitLength))
			})
		})
	})

	Context("Raft Node Election", func() {
		Context("RaftNode LeaderHeartbeatMonitor Timeouts", func() {
			node := &raft.RaftNode{}
			term_0 := int32(0)
			var setupVars MockSetupVars

			When("Raft Node doesn't receive leader heartbeat for the leader heartbeat duration", func() {

				BeforeEach(func() {

				})

				AfterEach(func() {
					defer setupVars.Ctrls.FileCtrl.Finish()
					defer setupVars.Ctrls.ElectionCtrl.Finish()
					defer setupVars.Ctrls.HeartCtrl.Finish()
					defer setupVars.Ctrls.RpcCtrl.Finish()
					defer node.Close()
				})

				It("Should become a candidate", func() {
					setupVars = SetupLeaderHeartbeatTimeout()
					node = setupVars.Node
					term_0 = setupVars.Term_0
					CheckIfEventTriggered(node, raft.BecameCandidate, raft.RaftEventDetails{CurrentTerm: term_0 + 1})
					Succeed()
				})

				It("Should Vote for itself", func() {
					setupVars = SetupLeaderHeartbeatTimeout()
					node = setupVars.Node
					CheckIfEventTriggered(node, raft.BecameCandidate, raft.RaftEventDetails{VotedFor: node.Id})
				})

				It("Should Increment the current term", func() {
					setupVars = SetupLeaderHeartbeatTimeout()
					node = setupVars.Node
					term_0 = setupVars.Term_0
					for {
						if node.CurrentTerm == term_0+1 {
							break
						}
					}
					Succeed()
				})

				It("Should Request votes from peers", func() {
					setupVars = SetupLeaderHeartbeatTimeout()
					node = setupVars.Node
					term_0 = setupVars.Term_0
					for {
						if len(node.ElectionMgr.GetReceivedVotes()) == len(node.Peers) {
							break
						}
					}
					Succeed()
				})
			})
		})

		Context("Candidate starts an election", func() {
			Context("Raft Node Election Overtimes", func() {
				node := &raft.RaftNode{}
				term_0 := int32(0)
				var setupVars MockSetupVars
				When("Candidate is not able to reach to a conclusion within the election allowed time", func() {
					BeforeEach(func() {

					})
					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						defer node.Close()
					})

					It("Should restart election", func() {
						setupVars = SetupPeerTakesTooMuchTimeToRespond()
						node = setupVars.Node
						term_0 = setupVars.Term_0

						CheckIfEventTriggered(node, raft.ElectionRestarted, raft.RaftEventDetails{})
						CheckIfEventTriggered(node, raft.ElectionTimerStarted, raft.RaftEventDetails{CurrentTerm: term_0 + 2})

					})

				})
			})

			Context("Collecting votes", func() {
				var setupVars MockSetupVars
				node := &raft.RaftNode{}
				When("Majority votes in favor", func() {
					BeforeEach(func() {

					})

					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						defer node.Close()
					})

					It("should become a leader", func() {
						setupVars = SetupMajorityVotesInFavor()
						node = setupVars.Node
						CheckIfEventTriggered(node, raft.BecameLeader, raft.RaftEventDetails{})
						Succeed()
					})

					It("Should cancel the election timer", func() {
						setupVars = SetupLeaderSendsHeartbeatsOnElectionConclusion()
						node = setupVars.Node

						CheckIfEventTriggered(node, raft.ElectionTimerStopped, raft.RaftEventDetails{})
						CheckIfEventTriggered(node, raft.LogRequestSent, raft.RaftEventDetails{Peer: node.Peers[0].Id})
						CheckIfEventTriggered(node, raft.LogRequestSent, raft.RaftEventDetails{Peer: node.Peers[1].Id})
						Succeed()

					})

					It("should replicate logs to all its peers", func() {
						setupVars = SetupLeaderSendsHeartbeatsOnElectionConclusion()
						node = setupVars.Node
						CheckIfEventTriggered(node, raft.LogRequestSent, raft.RaftEventDetails{Peer: node.Peers[0].Id})
						CheckIfEventTriggered(node, raft.LogRequestSent, raft.RaftEventDetails{Peer: node.Peers[1].Id})
					})
				})

				When("The candidate receives a vote response with a higher term than its current term", func() {

					BeforeEach(func() {

					})
					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						defer node.Close()
					})

					It("Should become a follower", func() {
						setupVars = SetupCandidateReceivesVoteResponseWithHigherTerm()
						node = setupVars.Node
						CheckIfEventTriggered(node, raft.ElectionTimerStopped, raft.RaftEventDetails{})

						Expect(node.CurrentRole).To(Equal(raft.FOLLOWER))
					})

					It("Should Make its current term equal to the one in the voteResponse", func() {
						setupVars = SetupCandidateReceivesVoteResponseWithHigherTerm()
						node = setupVars.Node
						CheckIfEventTriggered(node, raft.BecameFollower, raft.RaftEventDetails{CurrentTerm: int32(10)})
					})

					It("Should cancel election timer", func() {

						setupVars = SetupCandidateReceivesVoteResponseWithHigherTerm()
						node = setupVars.Node

						CheckIfEventTriggered(node, raft.ElectionTimerStopped, raft.RaftEventDetails{})

						Succeed()
					})
				})
			})

			Context("Peers receiving vote requests", func() {
				var setupVars MockSetupVars
				node := &raft.RaftNode{}

				When("The candidate's log and term are both ok", func() {

					BeforeEach(func() {

					})
					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						defer node.Close()
					})

					It("Should vote yes", func() {
						setupVars = SetupPeerReceivingCandidateVoteRequest()
						node = setupVars.Node

						node.ElectionMgr.GetResponseForVoteRequest(node, &pb.VoteRequest{
							CurrentTerm: 10,
							NodeId:      "",
							LogLength:   0,
							LastTerm:    0,
						})

						CheckIfEventTriggered(node, raft.VoteGranted, raft.RaftEventDetails{})

						Expect(node.VotedFor).To(Equal(setupVars.LeaderId))
					})
				})

				When("The candidate's log and term are not ok", func() {
					BeforeEach(func() {

					})
					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						defer node.Close()
					})

					It("Should vote no", func() {
						setupVars = SetupPeerReceivingCandidateVoteRequest()
						node = setupVars.Node

						node.ElectionMgr.GetResponseForVoteRequest(node, &pb.VoteRequest{
							CurrentTerm: 0,
							NodeId:      "1",
							LogLength:   0,
							LastTerm:    10,
						})

						CheckIfEventTriggered(node, raft.VoteNotGranted, raft.RaftEventDetails{})

						Expect(node.VotedFor).To(Equal(""))
					})
				})
			})

			Context("Candidate falls back to being a follower", func() {
				var setupVars MockSetupVars
				node := &raft.RaftNode{}
				When("A legitimate leader sends a heartbeat", func() {
					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						defer node.Close()
					})
					It("Should become a follower", func() {
						setupVars = SetupMajorityVotesAgainst()
						node = setupVars.Node

						CheckIfEventTriggered(node, raft.ElectionTimerStarted, raft.RaftEventDetails{})
						node.ElectionMgr.GetLeaderHeartChannel() <- &raft.RaftNode{
							Node: raft.Node{
								CurrentTerm: int32(10),
							},
						}

						CheckIfEventTriggered(node, raft.BecameFollower, raft.RaftEventDetails{CurrentTerm: int32(10)})
					})
				})
			})

			When("Candidate receives vote request from another candidate", func() {

				var setupVars MockSetupVars
				node := &raft.RaftNode{}

				BeforeEach(func() {

				})
				AfterEach(func() {
					setupVars.Ctrls.FileCtrl.Finish()
					setupVars.Ctrls.ElectionCtrl.Finish()
					setupVars.Ctrls.HeartCtrl.Finish()
					setupVars.Ctrls.RpcCtrl.Finish()
					defer node.Close()
				})

				It("should reject the vote request", func() {
					setupVars = SetupCandidateRequestsVoteFromCandidate()
					node = setupVars.Node
					node.ElectionMgr.BecomeACandidate(node)
					go func() {
						<-node.ElectionDone
					}()
					voteRequest := &pb.VoteRequest{
						NodeId:      "2",
						CurrentTerm: 9999,
						LogLength:   9999,
						LastTerm:    9998,
					}
					voteResponse, _ := node.ElectionMgr.GetResponseForVoteRequest(node, voteRequest)
					Expect(voteResponse.VoteGranted).To(Equal(false))
				})
			})
		})
	})

	Context("Broadcasting Messages", func() {
		var setupVars MockSetupVars
		node := &raft.RaftNode{}
		When("A Leader receives a broadcast request", func() {

			BeforeEach(func() {

			})
			AfterEach(func() {
				setupVars.Ctrls.FileCtrl.Finish()
				setupVars.Ctrls.ElectionCtrl.Finish()
				setupVars.Ctrls.HeartCtrl.Finish()
				setupVars.Ctrls.RpcCtrl.Finish()
				defer node.Close()
			})

			It("Should Append the entry to its log", func() {

				setupVars = SetupLeaderReceivesABroadcastRequest()
				node = setupVars.Node

				msg := GetBroadcastMsg()

				go func() {
					<-node.ElectionDone
				}()

				go func() {
					<-node.HeartDone
				}()

				Expect(len(node.Logs)).To(Equal(0))
				node.LogMgr.RespondToBroadcastMsgRequest(node, msg)
				Expect(len(node.Logs)).To(Equal(1))
			})

			It("Should update the Acknowledged log length for itself", func() {

				setupVars = SetupLeaderReceivesABroadcastRequest()
				node = setupVars.Node

				msg := GetBroadcastMsg()

				go func() {
					<-node.ElectionDone
				}()

				go func() {
					<-node.HeartDone
				}()

				Expect(len(node.Logs)).To(Equal(0))
				node.LogMgr.RespondToBroadcastMsgRequest(node, msg)
				Expect(node.AckedLength[node.Id]).To(Equal(int32(1)))
			})

			It("Should Replicate the log to its followers after appending to its own logs", func() {

				setupVars = SetupLeaderReceivesABroadcastRequest()
				node = setupVars.Node

				msg := GetBroadcastMsg()

				go func() {
					<-node.ElectionDone
				}()

				go func() {
					<-node.HeartDone
				}()

				Expect(len(node.Logs)).To(Equal(0))
				node.LogMgr.RespondToBroadcastMsgRequest(node, msg)

				CheckIfEventTriggered(node, raft.LogRequestSent, raft.RaftEventDetails{Peer: node.Peers[0].Id})
				CheckIfEventTriggered(node, raft.LogRequestSent, raft.RaftEventDetails{Peer: node.Peers[1].Id})
			})
		})

		When("A Follower receives a broadcast request", func() {
			BeforeEach(func() {

			})
			It("Should forward the request to the leader node", func() {
				setupVars = SetupFollowerReceivesBroadcastRequest()
				node = setupVars.Node

				msg := GetBroadcastMsg()

				go func() {
					<-node.ElectionDone
				}()

				go func() {
					<-node.HeartDone
				}()

				Expect(len(node.Logs)).To(Equal(0))
				node.LogMgr.RespondToBroadcastMsgRequest(node, msg)

				CheckIfEventTriggered(node, raft.ForwardedBroadcastReq, raft.RaftEventDetails{})
			})
		})
	})

	Context("Log Replication", func() {
		var setupVars MockSetupVars
		node := &raft.RaftNode{}
		When("Leader is replicating log to followers", func() {

			BeforeEach(func() {

			})
			AfterEach(func() {
				setupVars.Ctrls.FileCtrl.Finish()
				setupVars.Ctrls.ElectionCtrl.Finish()
				setupVars.Ctrls.HeartCtrl.Finish()
				setupVars.Ctrls.RpcCtrl.Finish()
				defer node.Close()
			})

			Specify(`The log request should have leaderId, 
			currentTerm, 
			index of the last sent log entry,
			term of the previous log entry,
			commitLength and
			suffix of log entries`, func() {

				setupVars = SetupMajorityVotesInFavor()
				node = setupVars.Node

				CheckIfEventTriggered(node, raft.LogRequestSent, raft.RaftEventDetails{Peer: node.Peers[0].Id})
				CheckIfEventTriggered(node, raft.LogRequestSent, raft.RaftEventDetails{Peer: node.Peers[1].Id})

				Expect((*setupVars.SentLogReplicationReq).CurrentTerm).To(Equal(node.CurrentTerm))
				Succeed()
			})
		})

		When("Follower receives a log replication request", func() {
			var setupVars MockSetupVars
			node := &raft.RaftNode{}
			When("the received log is ok", func() {

				BeforeEach(func() {

				})
				AfterEach(func() {
					setupVars.Ctrls.FileCtrl.Finish()
					setupVars.Ctrls.ElectionCtrl.Finish()
					setupVars.Ctrls.HeartCtrl.Finish()
					setupVars.Ctrls.RpcCtrl.Finish()
					defer node.Close()
				})

				It("Should update its currentTerm to to new term", func() {
					setupVars = SetupFollowerReceivesLogReplicationRequest()
					logReplReq := *setupVars.SentLogReplicationReq
					node = setupVars.Node

					node.LogMgr.RespondToLogReplicationRequest(node, logReplReq)

					Expect(node.CurrentTerm).To(Equal(logReplReq.CurrentTerm))
					Expect(node.CurrentRole).To(Equal(raft.FOLLOWER))
				})

				It("Should cancel the election it started (if any) if the leader is genuine", func() {
					setupVars = SetupFollowerReceivesLogReplicationRequest()
					logReplReq := *setupVars.SentLogReplicationReq
					node = setupVars.Node
					node.ElectionMgr.BecomeACandidate(node)

					logReplReq.CurrentTerm = node.CurrentTerm + 1
					Expect(node.ElectionInProgress).To(Equal(true))
					Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))

					go func() {
						node.ElectionMgr.ManageElection(node)
					}()

					node.LogMgr.RespondToLogReplicationRequest(node, logReplReq)

					CheckIfEventTriggered(node, raft.BecameFollower, raft.RaftEventDetails{CurrentTerm: logReplReq.CurrentTerm})
				})

				It("Should append the entry to its log", func() {
					setupVars = SetupFollowerReceivesLogReplicationRequest()
					logReplReq := *setupVars.SentLogReplicationReq
					node = setupVars.Node

					Expect(len(node.Logs)).To(Equal(0))

					node.CurrentTerm = logReplReq.GetCurrentTerm()
					node.LogMgr.RespondToLogReplicationRequest(node, logReplReq)

					Expect(len(node.Logs)).To(Equal(2))
				})

				It("Should update its Commit length of the log", func() {
					setupVars = SetupFollowerReceivesLogReplicationRequest()
					logReplReq := *setupVars.SentLogReplicationReq
					node = setupVars.Node

					Expect(len(node.Logs)).To(Equal(0))

					node.CurrentTerm = logReplReq.GetCurrentTerm()
					node.LogMgr.RespondToLogReplicationRequest(node, logReplReq)

					Expect(len(node.Logs)).To(Equal(2))
					Expect(node.CommitLength).To(Equal(int32(2)))
				})

				It("Should send the log response back to the leader with the acknowledged new length", func() {
					setupVars = SetupFollowerReceivesLogReplicationRequest()
					logReplReq := *setupVars.SentLogReplicationReq
					node = setupVars.Node

					Expect(len(node.Logs)).To(Equal(0))

					node.CurrentTerm = logReplReq.GetCurrentTerm()
					logResponse, _ := node.LogMgr.RespondToLogReplicationRequest(node, logReplReq)

					Expect(logResponse.AckLength).To(Equal(int32(2)))
				})

			})

			When("The received log or term is not ok", func() {

				BeforeEach(func() {

				})
				var setupVars MockSetupVars
				node := &raft.RaftNode{}
				AfterEach(func() {
					setupVars.Ctrls.FileCtrl.Finish()
					setupVars.Ctrls.ElectionCtrl.Finish()
					setupVars.Ctrls.HeartCtrl.Finish()
					setupVars.Ctrls.RpcCtrl.Finish()
					defer node.Close()
				})

				It(`Should send its currentTerm
			acknowledged length as 0
			acknowledgment as false to the leader`, func() {

					setupVars = SetupLeaderLogsAreNotOK()
					logReplReq := *setupVars.SentLogReplicationReq
					node = setupVars.Node

					node.CurrentTerm = logReplReq.GetCurrentTerm()
					logResponse, _ := node.LogMgr.RespondToLogReplicationRequest(node, logReplReq)

					Expect(logResponse.AckLength).To(Equal(int32(0)))
					Expect(logResponse.Success).To(Equal(false))
				})
			})

		})
	})

	Context("A node updates its log based on the log received from leader", func() {
		//TODO
	})

	Context("Leader Receives log acknowledgments", func() {
		When("The term in the log acknowledgment is more than leader's currentTerm", func() {})
		When("The term in the log acknowledgment is ok and log replication has been acknowledged by the follower", func() {

			When("log replication is acknowledged by the follower", func() {

				BeforeEach(func() {

				})
				var setupVars MockSetupVars
				node := &raft.RaftNode{}

				AfterEach(func() {
					setupVars.Ctrls.FileCtrl.Finish()
					setupVars.Ctrls.HeartCtrl.Finish()
					setupVars.Ctrls.RpcCtrl.Finish()
					defer node.Close()
				})

				It(`Should update its acknowledged and sentLength for the follower`, func() {

					setupVars = SetupLeaderReceivingLogReplicationAck()
					receivedLogResponses := setupVars.ReceivedLogResponse
					node = setupVars.Node
					node.LogMgr.ReplicateLog(node, node.Peers[0])
					node.LogMgr.ReplicateLog(node, node.Peers[1])

					go func() bool {
						done := <-node.HeartDone
						return done
					}()

					Expect(node.AckedLength[node.Peers[0].Id]).To(Equal((*receivedLogResponses)[node.Peers[0].Id].AckLength))
					Expect(node.AckedLength[node.Peers[1].Id]).To(Equal((*receivedLogResponses)[node.Peers[1].Id].AckLength))

					Expect(node.SentLength[node.Peers[0].Id]).To(Equal((*receivedLogResponses)[node.Peers[0].Id].AckLength))
					Expect(node.SentLength[node.Peers[1].Id]).To(Equal((*receivedLogResponses)[node.Peers[1].Id].AckLength))

				})

				It("should commit the logs", func() {

					setupVars = SetupLeaderReceivingLogReplicationAck()
					node = setupVars.Node
					node.LogMgr.ReplicateLog(node, node.Peers[0])
					node.LogMgr.ReplicateLog(node, node.Peers[1])

					go func() bool {
						done := <-node.HeartDone
						return done
					}()

					Expect(node.CommitLength).To(Equal(int32(len(node.Logs))))
				})
			})

			When("log replication is not acknowledged by the follower", func() {

				BeforeEach(func() {

				})
				var setupVars MockSetupVars
				node := &raft.RaftNode{}

				AfterEach(func() {
					setupVars.Ctrls.FileCtrl.Finish()
					setupVars.Ctrls.HeartCtrl.Finish()
					setupVars.Ctrls.RpcCtrl.Finish()
					defer node.Close()
				})

				It("Should update the sent length for the follower to one less than the previously attempted sent length value of the log", func() {

					setupVars = SetupLogReplicationNotAckByFollower()
					node = setupVars.Node
					node.SentLength[node.Peers[0].Id] = int32(2)
					node.SentLength[node.Peers[1].Id] = int32(2)
					node.LogMgr.ReplicateLog(node, node.Peers[0])
					node.LogMgr.ReplicateLog(node, node.Peers[1])

					go func() bool {
						done := <-node.HeartDone
						return done
					}()

					for {
						if node.SentLength[node.Peers[0].Id] == int32(1) &&
							node.SentLength[node.Peers[1].Id] == int32(1) {
							break
						}
						time.Sleep(200 * time.Microsecond)
					}
					Succeed()
				})

			})

		})

	})

})
