package raftelection

import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xaraphix/Sif/internal/raft"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
)

type ElectionManager struct {
	ElectionTimeoutDuration time.Duration
	ElectionTimerOff        bool

	VotesReceived []*pb.VoteResponse

	electionTimerDone       chan bool
	votesResponse           chan *pb.VoteResponse
	conclusionDone          chan bool
	askForVotesDone         chan bool
	electionTimedOut        chan bool
	followerAccordingToPeer chan *pb.VoteResponse
	leader                  chan bool
	follower                chan bool
	leaderHeartbeatChannel  chan *raft.RaftNode
}

func NewElectionManager() raft.RaftElection {
	return &ElectionManager{
		ElectionTimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
		VotesReceived:           []*pb.VoteResponse{},
		leaderHeartbeatChannel:  make(chan *raft.RaftNode),
	}
}

func (em *ElectionManager) SetElectionTimeoutDuration(timeoutDuration time.Duration) {
	em.ElectionTimeoutDuration = timeoutDuration
}

func (em *ElectionManager) StartElection(rn *raft.RaftNode) {
	em.initChannels()
	em.BecomeACandidate(rn)
	em.askForVotes(rn)
	em.startElectionTimer(rn)
	em.concludeFromReceivedVotes(rn)
	restartElection := em.ManageElection(rn)
	em.wrapUpElection(rn, restartElection)
}

func (em *ElectionManager) wrapUpElection(rn *raft.RaftNode, restartElection bool) {
	if restartElection {
		raft.LogEvent(raft.ElectionRestarted, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
		em.StartElection(rn)
	} else {
		rn.ElectionInProgress = false
	}
}

func (em *ElectionManager) initChannels() {
	em.votesResponse = make(chan *pb.VoteResponse)
	em.electionTimedOut = make(chan bool)
	em.followerAccordingToPeer = make(chan *pb.VoteResponse)
	em.leader = make(chan bool)
	em.follower = make(chan bool)
}

func (em *ElectionManager) BecomeACandidate(rn *raft.RaftNode) {
	rn.ElectionInProgress = true
	rn.CurrentRole = raft.CANDIDATE
	rn.VotedFor = rn.Id
	rn.CurrentTerm = rn.CurrentTerm + 1
	em.VotesReceived = nil

  raft.LogEvent(raft.BecameCandidate, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, VotedFor: rn.Id, CurrentRole: rn.CurrentRole})
	logrus.WithFields(logrus.Fields{
		"Started By":    rn.Config.InstanceName(),
		"Started By Id": rn.Config.InstanceId(),
		"Term":          rn.CurrentTerm,
	}).Debug("Starting election")
}

func (em *ElectionManager) ManageElection(rn *raft.RaftNode) bool {
	for {
		select {
		case vr := <-em.followerAccordingToPeer:
			em.becomeAFollowerAccordingToPeer(rn, vr)
			return false
		case <-em.leader:
			em.becomeALeader(rn)
			return false
		case l := <-em.leaderHeartbeatChannel:
			em.becomeAFollowerAccordingToLeader(rn, l)
			return false
		case <-em.follower:
			em.becomeAFollower(rn)
			return false
		case <-em.electionTimedOut:
			raft.LogEvent(raft.ElectionTimerStopped, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
			return true
		}
	}
}

func (em *ElectionManager) becomeALeader(rn *raft.RaftNode) {
	if rn.ElectionInProgress == true {
		rn.CurrentRole = raft.LEADER
		raft.LogEvent(raft.BecameLeader, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
		em.replicateLogs(rn)
		rn.Heart.StartBeating(rn)
		rn.ElectionInProgress = false
		raft.LogEvent(raft.ElectionTimerStopped, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})

		logrus.WithFields(logrus.Fields{
			"Name": rn.Config.InstanceName(),
			"Id":   rn.Config.InstanceId(),
		}).Debug("I became the leader")
	}
}

func (em *ElectionManager) becomeAFollower(rn *raft.RaftNode) {
	if rn.ElectionInProgress == true {
		rn.CurrentRole = raft.FOLLOWER
		rn.ElectionInProgress = false
		rn.VotedFor = ""
		raft.LogEvent(raft.BecameFollower, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
		raft.LogEvent(raft.ElectionTimerStopped, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
		logrus.WithFields(logrus.Fields{
			"Name": rn.Config.InstanceName(),
		}).Debug("I became a follower because peers voted against me")
	}
}

func (em *ElectionManager) replicateLogs(rn *raft.RaftNode) {
	go func(n *raft.RaftNode) {
		for _, peer := range rn.Peers {
			rn.LogMgr.ReplicateLog(rn, peer)
		}
	}(rn)
}

func (em *ElectionManager) becomeAFollowerAccordingToLeader(rn *raft.RaftNode, leader *raft.RaftNode) {
	if rn.ElectionInProgress == true {
		rn.CurrentRole = raft.FOLLOWER
		rn.CurrentTerm = leader.CurrentTerm
    rn.CurrentLeader = leader.Id
		rn.VotedFor = ""
		rn.ElectionInProgress = false
    raft.LogEvent(raft.BecameFollower, raft.RaftEventDetails{Id: rn.Id, CurrentLeader: leader.Id, CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
		raft.LogEvent(raft.ElectionTimerStopped, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
		logrus.WithFields(logrus.Fields{
			"Name": rn.Config.InstanceName(),
		}).Debug("I became a follower according to leader")
	}
}

func (em *ElectionManager) becomeAFollowerAccordingToPeer(rn *raft.RaftNode, v *pb.VoteResponse) {
	if rn.ElectionInProgress == true {
		rn.CurrentTerm = v.Term
		rn.CurrentRole = raft.FOLLOWER
		rn.VotedFor = ""
		rn.ElectionInProgress = false
		raft.LogEvent(raft.BecameFollower, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
		raft.LogEvent(raft.ElectionTimerStopped, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
		logrus.WithFields(logrus.Fields{
			"Name": rn.Config.InstanceName(),
		}).Debug("I became a follower according to peer")
	}
}

func (em *ElectionManager) startElectionTimer(rn *raft.RaftNode) {
	em.electionTimedOut = make(chan bool)
	go func() {
		raft.LogEvent(raft.ElectionTimerStarted, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
		timer := time.NewTimer(em.ElectionTimeoutDuration)
		for {
			select {
			case <-timer.C:
				if rn.ElectionInProgress {
					raft.LogEvent(raft.ElectionTimerStopped, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
					timer.Stop()
					em.electionTimedOut <- true
					return
				} else {
					raft.LogEvent(raft.ElectionTimerStopped, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole})
					timer.Stop()
					return
				}
			case <-rn.ElectionDone:
				rn.ElectionInProgress = false
				break
			}
		}
	}()
}

func (em *ElectionManager) askForVotes(rn *raft.RaftNode) {
	voteRequest := em.GenerateVoteRequest(rn)
	for idx, peer := range rn.Peers {
		go func(i int, p raft.Peer, n *raft.RaftNode) {
			voteResponse := n.RPCAdapter.RequestVoteFromPeer(p, voteRequest)
			em.votesResponse <- voteResponse
		}(idx, peer, rn)
	}
}

func (em *ElectionManager) GetReceivedVotes() []*pb.VoteResponse {
	return em.VotesReceived
}

func (em *ElectionManager) GetLeaderHeartChannel() chan *raft.RaftNode {
	return em.leaderHeartbeatChannel
}

func (em *ElectionManager) GetResponseForVoteRequest(rn *raft.RaftNode, vr *pb.VoteRequest) (*pb.VoteResponse, error) {
  raft.LogEvent(raft.VoteRequestReceived, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole, Peer: vr.NodeId})
	voteResponse := em.getVoteResponseForVoteRequest(rn, vr)

	if voteResponse.VoteGranted {
    raft.LogEvent(raft.VoteGranted, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole, Peer: vr.NodeId})
	} else {
    raft.LogEvent(raft.VoteNotGranted, raft.RaftEventDetails{CurrentTerm: rn.CurrentTerm, CurrentRole: rn.CurrentRole, Peer: vr.NodeId})
	}
	return voteResponse, nil
}

func (em *ElectionManager) GetElectionTimeoutDuration() time.Duration {
	return em.ElectionTimeoutDuration
}
