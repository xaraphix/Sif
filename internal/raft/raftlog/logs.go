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
		prevLogTerm = rn.Logs[i].Term
	}

	replicateLogsRequest := raft.LogRequest{
		LeaderId:     rn.Id,
		CurrentTerm:  rn.CurrentTerm,
		SentLength:   i,
		PrevLogTerm:  prevLogTerm,
		CommitLength: rn.CommitLength,
		Entries:      &entries,
	}

	rn.RPCAdapter.ReplicateLog(peer, replicateLogsRequest)
	rn.SendSignal(raft.LogRequestSent)
}

func (l *LogMgr) BroadcastMessage(rn *raft.RaftNode, msg map[string]interface{}) {

	if rn.CurrentRole == raft.LEADER {
		rn.SendSignal(raft.MsgAppendedToLogs)
		rn.Logs = append(rn.Logs, raft.Log{
			Term: rn.CurrentTerm,
			Message: msg,
		})

		rn.AckedLength[rn.Id] = int32(len(rn.Logs))

		for _, peer := range rn.Peers {
			go func (n *raft.RaftNode, p raft.Peer)  {
				l.ReplicateLog(n, p)
			}(rn, peer)
		}
	} else {
		//raftRPCAdapter forward request
	}
}
