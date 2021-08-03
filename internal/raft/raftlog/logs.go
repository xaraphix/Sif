package raftlog

import "github.com/xaraphix/Sif/internal/raft"

type LogMgr struct {
		
}

func (l *LogMgr) GetLogs() []raft.Log {
	return nil
}

func (l *LogMgr) GetLog(idx int32) raft.Log {
	return raft.Log{}
}

func (l *LogMgr) ReplicateLog(rn *raft.RaftNode, peer raft.Peer) {

}
