package raftlog

import (
	"github.com/xaraphix/Sif/internal/raft"
)

type LogMgr struct {
}

func (l *LogMgr) GetLogs() []raft.Log {
	return nil
}

func (l *LogMgr) GetLog(rn *raft.RaftNode, idx int32) raft.Log {
	return raft.Log{}
}

func (l *LogMgr) ReplicateLog(rn *raft.RaftNode, peer raft.Peer) {
	i := rn.SentLength[peer.Id]
	entries := rn.Logs[i:len(rn.Logs)]
	prevLogTerm := int32(0)

	if i > 0 {
		prevLogTerm = rn.Logs[i-1].Term
	}

	replicateLogsRequest := raft.LogRequest{
		LeaderId:     rn.Id,
		Term:         rn.CurrentTerm,
		LogLength:    i,
		LogTerm:      prevLogTerm,
		CommitLength: rn.CommitLength,
		Entries:      &entries,
	}

	rn.RPCAdapter.ReplicateLog(peer, replicateLogsRequest)
	rn.SendSignal(raft.LogRequestSent)
}

func (l *LogMgr) RespondToBroadcastMsgRequest(rn *raft.RaftNode, msg map[string]interface{}) raft.BroadcastMessageResponse {
	if rn.CurrentRole == raft.LEADER {
		rn.SendSignal(raft.MsgAppendedToLogs)
		rn.Logs = append(rn.Logs, raft.Log{
			Term:    rn.CurrentTerm,
			Message: msg,
		})

		rn.AckedLength[rn.Id] = int32(len(rn.Logs))

		for _, peer := range rn.Peers {
			go func(n *raft.RaftNode, p raft.Peer) {
				l.ReplicateLog(n, p)
			}(rn, peer)
		}

		return raft.BroadcastMessageResponse{}
	} else {
		leaderPeer := rn.GetPeerById(rn.CurrentLeader)
		return rn.RPCAdapter.BroadcastMessage(leaderPeer, msg)
	}
}

func (l *LogMgr) RespondToLogReplicationRequest(rn *raft.RaftNode, lr raft.LogRequest) raft.LogResponse {
	if lr.Term > rn.CurrentTerm {
		rn.CurrentTerm = lr.Term
		rn.VotedFor = 0
	}

	logOk := int32(len(rn.Logs)) >= lr.LogLength

	if logOk && lr.LogLength > 0 {
		logOk = lr.LogTerm == rn.Logs[lr.LogLength-1].Term
	}

	if rn.ElectionInProgress && rn.CurrentTerm < lr.Term && logOk {
		rn.ElectionMgr.GetLeaderHeartChannel() <- raft.RaftNode{
			Node: raft.Node{
				Id:          lr.LeaderId,
				CurrentTerm: lr.Term,
			},
		}
		rn.CurrentTerm = lr.Term
	}

	if lr.Term == rn.CurrentTerm && logOk {
		rn.CurrentRole = raft.FOLLOWER
		rn.CurrentLeader = lr.LeaderId
		//TODO l.appendEntries()
		ack := lr.LogLength + int32(len(*lr.Entries))
		return raft.LogResponse{
			FollowerId: rn.Id,
			Term:       rn.CurrentTerm,
			AckLength:  ack,
			Success:    true,
		}
	} else {
		return raft.LogResponse{
			FollowerId: rn.Id,
			Term:       rn.CurrentTerm,
			AckLength:  0,
			Success:    false,
		}
	}
}

func (l *LogMgr) ReceiveLogAcknowledgements(rn *raft.RaftNode, lr raft.LogResponse) {
	if lr.Term == rn.CurrentTerm && rn.CurrentRole == raft.LEADER {
		if lr.Success {
			rn.SentLength[lr.FollowerId] = lr.AckLength
			rn.AckedLength[lr.FollowerId] = lr.AckLength
			l.commitLogEntries(rn)
		} else if rn.SentLength[lr.FollowerId] > 0 {
			rn.SentLength[lr.FollowerId] = rn.SentLength[lr.FollowerId] - 1
			follower := rn.GetPeerById(lr.FollowerId)
			l.ReplicateLog(rn, follower)
		}
	} else if lr.Term > rn.CurrentTerm {
		rn.CurrentTerm = lr.Term
		rn.CurrentRole = raft.FOLLOWER
		rn.VotedFor = 0
	}
}

func (l *LogMgr) commitLogEntries(rn *raft.RaftNode) {
	// minAcks := math.Ceil((len(rn.Peers) + 1) /2)
	//TODO
}

//GIVEN an ackLength return the peers who have atleast acked till that length
func countOfNodesWithAckLengthGTE(rn *raft.RaftNode, ackLength int32) int {
	count := 0
	for _, peer := range rn.Peers {
		if rn.AckedLength[peer.Id] >= ackLength {
			count++
		}
	}

	return count
}

func (l *LogMgr) appendEntries(rn *raft.RaftNode, logLength int32, leaderCommitLength int32, entries []raft.Log) {
	if len(entries) > 0 && int32(len(rn.Logs)) > logLength {
		if rn.Logs[logLength].Term != entries[0].Term {
			//truncate logs
			rn.Logs = rn.Logs[0:logLength]
		}
	}

	if logLength+int32(len(entries)) > int32(len(rn.Logs)) {
		for i := int32(len(rn.Logs)) - logLength; i < int32(len(entries)); i++ {
			rn.Logs = append(rn.Logs, entries[i])
		}
	}

	if leaderCommitLength > rn.CommitLength {
		for i := rn.CommitLength; i < leaderCommitLength; i++ {
			//deliver rn.Logs[i].msg to the app TODO
		}

		rn.CommitLength = leaderCommitLength
	}
}
