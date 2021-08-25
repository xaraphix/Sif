package raftlog

import (
	"math"

	"github.com/sirupsen/logrus"
	"github.com/xaraphix/Sif/internal/raft"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"google.golang.org/protobuf/types/known/structpb"
)

type LogMgr struct {
}

func NewLogManager() *LogMgr {
	return &LogMgr{}
}

func (l *LogMgr) GetLogs() []*pb.Log {
	return nil
}

func (l *LogMgr) GetLog(rn *raft.RaftNode, idx int32) *pb.Log {
	return &pb.Log{}
}

func (l *LogMgr) ReplicateLog(rn *raft.RaftNode, peer raft.Peer) {
	i := rn.SentLength[peer.Id]
	entries := (rn.Logs)[i:len(rn.Logs)]
	prevLogTerm := int32(0)

	if i > 0 {
		prevLogTerm = (rn.Logs)[i-1].Term
	}

	replicateLogsRequest := &pb.LogRequest{
		LeaderId:     rn.Id,
		CurrentTerm:  rn.CurrentTerm,
		SentLength:   i,
		PrevLogTerm:  prevLogTerm,
		CommitLength: rn.CommitLength,
		Entries:      entries,
	}

	logResponse := rn.RPCAdapter.ReplicateLog(peer, replicateLogsRequest)
	rn.LogEvent(raft.LogRequestSent, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole, Peer: peer.Id})
	l.processLogAcknowledgements(rn, logResponse)
}

func (l *LogMgr) RespondToBroadcastMsgRequest(rn *raft.RaftNode, msg *structpb.Struct) (*pb.BroadcastMessageResponse, error) {
	if rn.CurrentRole == raft.LEADER {
		rn.LogEvent(raft.MsgAppendedToLogs, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole, CurrentLeader: rn.CurrentLeader})
		rn.Logs = append(rn.Logs, &pb.Log{
			Term:    rn.CurrentTerm,
			Message: msg,
		})

		rn.AckedLength[rn.Id] = int32(len(rn.Logs))

		for _, peer := range rn.Peers {
			go func(n *raft.RaftNode, p raft.Peer) {
				l.ReplicateLog(n, p)
			}(rn, peer)
		}

		return &pb.BroadcastMessageResponse{}, nil
	} else {
		leaderPeer := rn.GetPeerById(rn.CurrentLeader)
    rn.LogEvent(raft.ForwardedBroadcastReq, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole, CurrentLeader: rn.CurrentLeader})
		return rn.RPCAdapter.BroadcastMessage(leaderPeer, msg), nil
	}
}

func (l *LogMgr) RespondToLogReplicationRequest(rn *raft.RaftNode, lr *pb.LogRequest) (*pb.LogResponse, error) {

	if lr.CurrentTerm > rn.CurrentTerm {
		rn.CurrentTerm = lr.CurrentTerm
		rn.VotedFor = "" 
	}

	logOk := int32(len(rn.Logs)) >=lr.SentLength

	if logOk && lr.SentLength > 0 {
		logOk = lr.PrevLogTerm == rn.Logs[lr.SentLength-1].Term
	}

	if rn.ElectionInProgress && rn.CurrentTerm <= lr.CurrentTerm && logOk {
		rn.ElectionMgr.GetLeaderHeartChannel() <- &raft.RaftNode{
			Node: raft.Node{
				Id:          lr.LeaderId,
				CurrentTerm: lr.CurrentTerm,
			},
		}
	}

	if lr.CurrentTerm == rn.CurrentTerm && logOk {

		rn.CurrentRole = raft.FOLLOWER
		rn.CurrentLeader = lr.LeaderId

		logrus.WithFields(logrus.Fields{
			"LeaderId":      lr.LeaderId,
			"FollowerId":    rn.Id,
			"MyTerm":        rn.CurrentTerm,
			"Leader Term":   lr.CurrentTerm,
			"CurrentLeader": rn.CurrentLeader,
		}).Debug("Received Heartbeat from leader")

		rn.LeaderHeartbeatMonitor.Reset()
		l.appendEntries(rn, lr.SentLength, lr.CommitLength, lr.Entries)
		ack := lr.SentLength + int32(len(lr.Entries))
		return &pb.LogResponse{
			FollowerId: rn.Id,
			Term:       rn.CurrentTerm,
			AckLength:  ack,
			Success:    true,
		}, nil
	} else {
		return &pb.LogResponse{
			FollowerId: rn.Id,
			Term:       rn.CurrentTerm,
			AckLength:  0,
			Success:    false,
		}, nil
	}
}

func (l *LogMgr) processLogAcknowledgements(rn *raft.RaftNode, lr *pb.LogResponse) {
	if lr == nil {
		return
	}

	if lr.Term == rn.CurrentTerm && rn.CurrentRole == raft.LEADER {
		rn.LogAckMu.Lock()
		defer rn.LogAckMu.Unlock()
		if lr.Success {
			rn.SentLength[lr.FollowerId] = lr.AckLength
			rn.AckedLength[lr.FollowerId] = lr.AckLength
			l.commitLogEntries(rn)
		} else if rn.SentLength[lr.FollowerId] > 0 {
			rn.SentLength[lr.FollowerId] = rn.SentLength[lr.FollowerId] - 1
			follower := rn.GetPeerById(lr.FollowerId)
			go func() {
				l.ReplicateLog(rn, follower)
			}()
		}
	} else if lr.Term > rn.CurrentTerm {
		rn.CurrentTerm = lr.Term
		rn.CurrentRole = raft.FOLLOWER
		rn.VotedFor = ""
	}
}

func (l *LogMgr) commitLogEntries(rn *raft.RaftNode) {
	minAcks := int(math.Ceil(float64((len(rn.Peers) + 1) / 2)))
	maxReady := 0
	for i := 0; i < len(rn.Logs); i++ {
		if countOfNodesWithAckLengthGTE(rn, i) > minAcks && i > maxReady {
			maxReady = i + 1
		}
	}

	if maxReady > 0 && maxReady > int(rn.CommitLength) &&
		rn.Logs[maxReady-1].Term == rn.CurrentTerm {
		for i := int(rn.CommitLength); i < maxReady; i++ {
			l.deliverToApplication(rn, rn.Logs[i].Message)
		}

		rn.CommitLength = int32(maxReady)
	}
}

//GIVEN an ackLength return the peers who have atleast acked till that length
func countOfNodesWithAckLengthGTE(rn *raft.RaftNode, ackLength int) int {
	count := 0
	for _, peer := range rn.Peers {
		if int(rn.AckedLength[peer.Id]) >= ackLength {
			count++
		}
	}

	return count
}

func (l *LogMgr) deliverToApplication(rn *raft.RaftNode, msg *structpb.Struct) {
	logrus.WithFields(logrus.Fields{
		"Delivered By": rn.Id,
	}).Debug("Delivering msgs to application")

	rn.LogEvent(raft.DeliveredToApplication, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole, CurrentLeader: rn.CurrentLeader})
}

func (l *LogMgr) appendEntries(rn *raft.RaftNode, logLength int32, leaderCommitLength int32, entries []*pb.Log) {
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
			l.deliverToApplication(rn, rn.Logs[i].Message)
		}

		rn.CommitLength = leaderCommitLength
	}
}
