package raftlog

import (
	"github.com/sirupsen/logrus"
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

	replicateLogsRequest := raft.ReplicateLogRequest{
		LeaderId: rn.Id,
		CurrentTerm: rn.CurrentTerm,
		SentLength: i,
		PrevLogTerm: prevLogTerm,
		CommitLength: rn.CommitLength,
		Entries: &entries,
	}
	
	logrus.Info("===================== SEnding HEART BEAT")
	rn.RPCAdapter.ReplicateLog(peer, replicateLogsRequest)
}
