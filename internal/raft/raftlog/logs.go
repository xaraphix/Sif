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

func (l *LogMgr) RespondToBroadcastMsgRequest(rn *raft.RaftNode, msg map[string]interface{}) {

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
	} else {
		//TODO raftRPCAdapter forward request
	}
}

func (l *LogMgr) RespondToLogRequest(rn *raft.RaftNode, lr raft.LogRequest) raft.LogResponse {

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
				Id: lr.LeaderId,
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
