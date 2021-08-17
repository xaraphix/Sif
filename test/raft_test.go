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

			It("should reset currentTerm", func() {
				Expect(node.CurrentTerm).To(Equal(int32(0)))
			})

			It("should reset logs", func() {
				Expect(len(node.Logs)).To(Equal(0))
			})

			It("should reset VotedFor", func() {
				Expect(node.VotedFor).To(Equal(int32(0)))
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
				raft.DestructRaftNode(node)
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
					setupVars = SetupLeaderHeartbeatTimeout()
					node = setupVars.Node
					term_0 = setupVars.Term_0
				})

				AfterEach(func() {
					defer setupVars.Ctrls.FileCtrl.Finish()
					defer setupVars.Ctrls.ElectionCtrl.Finish()
					defer setupVars.Ctrls.HeartCtrl.Finish()
					defer setupVars.Ctrls.RpcCtrl.Finish()
					defer raft.DestructRaftNode(node)
				})

				It("Should become a candidate", func() {
					for e := range node.GetRaftSignalsChan() {
						if e == raft.ElectionTimerStarted {
							break
						}
					}
					Expect(node.CurrentRole).To(Equal(raft.CANDIDATE))
				})

				It("Should Vote for itself", func() {
					for e := range node.GetRaftSignalsChan() {
						if e == raft.ElectionTimerStarted {
							break
						}
					}
					Expect(node.VotedFor).To(Equal(node.Id))
				})

				It("Should Increment the current term", func() {
					for e := range node.GetRaftSignalsChan() {
						if e == raft.ElectionTimerStarted {
							break
						}
					}
					Expect(node.CurrentTerm).To(BeNumerically("==", term_0+1))
				})

				It("Should Request votes from peers", func() {
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
					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("Should restart election", func() {
						setupVars = SetupPeerTakesTooMuchTimeToRespond()
						node = setupVars.Node
						term_0 = setupVars.Term_0

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionRestarted {
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStarted {
								break
							}
						}

						Expect(node.CurrentTerm).To(Equal(term_0 + 2))
					})

				})
			})

			Context("Collecting votes", func() {
				var setupVars MockSetupVars
				node := &raft.RaftNode{}
				When("Majority votes in favor", func() {

					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("should become a leader", func() {
						setupVars = SetupMajorityVotesInFavor()
						node = setupVars.Node
						for e := range node.GetRaftSignalsChan() {
							if e == raft.BecameLeader {
								break
							}
						}
						Succeed()
					})

					It("Should cancel the election timer", func() {
						setupVars = SetupLeaderSendsHeartbeatsOnElectionConclusion()
						node = setupVars.Node

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStarted {
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStopped {
								break
							}
						}

						counter := 0
						for e := range node.GetRaftSignalsChan() {
							if e == raft.LogRequestSent {
								counter = counter + 1
								if counter == 2 {
									break
								}
							}
						}

						time.Sleep(100 * time.Millisecond)

						Succeed()

					})

					It("should replicate logs to all its peers", func() {

						setupVars = SetupLeaderSendsHeartbeatsOnElectionConclusion()
						node = setupVars.Node

						for e := range setupVars.Node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStarted {
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.BecameLeader {
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.HeartbeatStarted {
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.LogRequestSent {
								Succeed()
								break
							}
						}

					})
				})

				When("The candidate receives a vote response with a higher term than its current term", func() {

					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("Should become a follower", func() {
						setupVars = SetupCandidateReceivesVoteResponseWithHigherTerm()
						node = setupVars.Node
						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStopped {
								break
							}
						}

						Expect(node.CurrentRole).To(Equal(raft.FOLLOWER))
					})

					It("Should Make its current term equal to the one in the voteResponse", func() {
						setupVars = SetupCandidateReceivesVoteResponseWithHigherTerm()
						node = setupVars.Node
						testConfig := LoadTestRaftConfig()
						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStopped {
								break
							}
						}
						Expect(node.CurrentTerm).To(Equal((*setupVars.ReceivedVoteResponse)[testConfig.Peers()[0].Id].Term))
					})

					It("Should cancel election timer", func() {

						setupVars = SetupCandidateReceivesVoteResponseWithHigherTerm()
						node = setupVars.Node

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStopped {
								break
							}
						}

						Succeed()
					})
				})
			})

			Context("Peers receiving vote requests", func() {
				var setupVars MockSetupVars
				node := &raft.RaftNode{}

				When("The candidate's log and term are both ok", func() {

					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("Should vote yes", func() {
						setupVars = SetupPeerReceivingCandidateVoteRequest()
						node = setupVars.Node

						node.ElectionMgr.GetResponseForVoteRequest(node, &pb.VoteRequest{
							CurrentTerm: 10,
							NodeId:      0,
							LogLength:   0,
							LastTerm:    0,
						})

						for e := range node.GetRaftSignalsChan() {
							if e == raft.VoteGranted {
								break
							}
						}

						Expect(node.VotedFor).To(Equal(setupVars.LeaderId))
					})
				})

				When("The candidate's log and term are not ok", func() {
					AfterEach(func() {
						setupVars.Ctrls.FileCtrl.Finish()
						setupVars.Ctrls.ElectionCtrl.Finish()
						setupVars.Ctrls.HeartCtrl.Finish()
						setupVars.Ctrls.RpcCtrl.Finish()
						raft.DestructRaftNode(node)
					})

					It("Should vote no", func() {
						setupVars = SetupPeerReceivingCandidateVoteRequest()
						node = setupVars.Node

						node.ElectionMgr.GetResponseForVoteRequest(node, &pb.VoteRequest{
							CurrentTerm: 0,
							NodeId:      1,
							LogLength:   0,
							LastTerm:    10,
						})

						for e := range node.GetRaftSignalsChan() {
							if e == raft.VoteNotGranted {
								break
							}
						}

						Expect(node.VotedFor).To(Equal(int32(0)))
						raft.DestructRaftNode(node)
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
						raft.DestructRaftNode(node)
					})
					It("Should become a follower", func() {
						setupVars = SetupMajorityVotesAgainst()
						node = setupVars.Node

						for e := range node.GetRaftSignalsChan() {
							if e == raft.ElectionTimerStarted {
								node.ElectionMgr.GetLeaderHeartChannel() <- raft.RaftNode{
									Node: raft.Node{
										CurrentTerm: int32(10),
									},
								}
								break
							}
						}

						for e := range node.GetRaftSignalsChan() {
							if e == raft.BecameFollower {
								break
							}
						}

						Expect(node.CurrentRole).To(Equal(raft.FOLLOWER))
						Expect(node.CurrentTerm).To(Equal(int32(10)))
					})
				})
			})

			When("Candidate receives vote request from another candidate", func() {
				XIt("should reject the vote request", func() {
					Fail("Not Yet Implemented")
				})
			})
		})
	})

	Context("Broadcasting Messages", func() {
		When("A Leader receives a broadcast request", func() {

			XIt("Should Append the entry to its log", func() {

			})

			XIt("Should update the Acknowledged log length for itself", func() {

			})

			XIt("Should Replicate the log to its followers after appending to its own logs", func() {

			})
		})

		When("A Raft node is a leader", func() {
			XIt("Should send heart beats to its followers", func() {

			})
		})

		When("A Follower receives a broadcast request", func() {
			XIt("Should forward the request to the leader node", func() {

			})
		})
	})

	Context("Log Replication", func() {
		var setupVars MockSetupVars
		node := &raft.RaftNode{}
		When("Leader is replicating log to followers", func() {

			AfterEach(func() {
				setupVars.Ctrls.FileCtrl.Finish()
				setupVars.Ctrls.ElectionCtrl.Finish()
				setupVars.Ctrls.HeartCtrl.Finish()
				setupVars.Ctrls.RpcCtrl.Finish()
				raft.DestructRaftNode(node)
			})

			Specify(`The log request should have leaderId, 
			currentTerm, 
			index of the last sent log entry,
			term of the previous log entry,
			commitLength and
			suffix of log entries`, func() {

				setupVars = SetupMajorityVotesInFavor()
				node = setupVars.Node

				for e := range node.GetRaftSignalsChan() {
					if e == raft.LogRequestSent {
						break
					}
				}

				Expect((*setupVars.SentLogReplicationReq).CurrentTerm).To(Equal(node.CurrentTerm))
				Succeed()
			})
		})

		When("Follower receives a log replication request", func() {
			When("the received log is ok", func() {

				XIt("Should update its currentTerm to to new term", func() {

				})

				XIt("Should continue being the follower and cancel the election it started (if any)", func() {

				})

				XIt("Should update its acknowledged length of the log", func() {

				})

				XIt("Should append the entry to its log", func() {

				})

				XIt("Should send the log response back to the leader with the acknowledged new length", func() {

				})

			})

			When("The received log or term is not ok", func() {

				XIt(`Should send its currentTerm
			acknowledged length as 0
			acknowledgment as false to the leader`, func() {

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

				XIt(`Should update its acknowledged and 
					sentLength for the follower`, func() {

				})

				XIt("Should commit the log entries to persistent storage", func() {

				})
			})

			When("log replication is not acknowledged by the follower", func() {

				XIt("Should update the sent length for the follower to one less than the previously attempted sent length value of the log", func() {

				})
				XIt("Should send a log replication request to the follower with log length = last request log length - 1", func() {

				})
			})

		})

	})

	Context("Commiting log entries to persistent storage", func() {
		//TODO
		When("The commit is successful it should send the log message to the client of Sif", func() {

			XIt("Don't know yet what to do", func() {
				Fail("Not yet Implemented")
			})
		})
	})

})