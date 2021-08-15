package raftelection

import (
	"math/rand"
	"time"

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
	leaderHeartbeatChannel  chan raft.RaftNode
	peerVoteChannel         chan raft.RaftNode
}

func NewElectionManager() raft.RaftElection {
	return &ElectionManager{
		ElectionTimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
		VotesReceived:           []*pb.VoteResponse{},
		leaderHeartbeatChannel:  make(chan raft.RaftNode),
	}

}

func (em *ElectionManager) StartElection(rn *raft.RaftNode) {
	em.initChannels()
	em.becomeACandidate(rn)
	em.askForVotes(rn)
	em.startElectionTimer(rn)
	em.concludeFromReceivedVotes(rn)
	restartElection := em.handleElection(rn)
	em.wrapUpElection(rn, restartElection)
}

func (em *ElectionManager) wrapUpElection(rn *raft.RaftNode, restartElection bool) {
	if restartElection {
		rn.SendSignal(raft.ElectionRestarted)
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

func (em *ElectionManager) becomeACandidate(rn *raft.RaftNode) {
	rn.ElectionInProgress = true
	rn.CurrentRole = raft.CANDIDATE
	rn.VotedFor = rn.Id
	rn.CurrentTerm = rn.CurrentTerm + 1
	em.VotesReceived = nil
}

func (em *ElectionManager) handleElection(rn *raft.RaftNode) bool {
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
		case p := <-em.peerVoteChannel:
			em.becomeAFollowerAccordingToCandidatePeer(rn, p)
			return false
		case <-em.follower:
			em.becomeAFollower(rn)
			return false
		case <-em.electionTimedOut:
			rn.SendSignal(raft.ElectionTimerStopped)
			return true
		}
	}
}

func (em *ElectionManager) becomeALeader(rn *raft.RaftNode) {
	if rn.ElectionInProgress == true {
		rn.CurrentRole = raft.LEADER
		rn.SendSignal(raft.BecameLeader)
		em.replicateLogs(rn)
		rn.Heart.StartBeating(rn)
		rn.ElectionInProgress = false
		rn.SendSignal(raft.ElectionTimerStopped)
	}
}

func (em *ElectionManager) becomeAFollower(rn *raft.RaftNode) {
	if rn.ElectionInProgress == true {
		rn.CurrentRole = raft.FOLLOWER
		rn.ElectionInProgress = false
		rn.SendSignal(raft.ElectionTimerStopped)
	}
}

func (em *ElectionManager) replicateLogs(rn *raft.RaftNode) {
	go func(n *raft.RaftNode) {
		for _, peer := range rn.Peers {
			rn.LogMgr.ReplicateLog(rn, peer)
		}
		rn.SendSignal(raft.LogRequestSent)
	}(rn)
}

func (em *ElectionManager) becomeAFollowerAccordingToLeader(rn *raft.RaftNode, leader raft.RaftNode) {
	if rn.ElectionInProgress == true {
		rn.CurrentRole = raft.FOLLOWER
		rn.SendSignal(raft.BecameFollower)
		rn.CurrentTerm = leader.CurrentTerm
		rn.VotedFor = leader.Id
		rn.ElectionInProgress = false
		rn.SendSignal(raft.ElectionTimerStopped)
	}
}

func (em *ElectionManager) becomeAFollowerAccordingToCandidatePeer(rn *raft.RaftNode, peer raft.RaftNode) {
	if rn.ElectionInProgress == true {
		rn.CurrentRole = raft.FOLLOWER
		rn.SendSignal(raft.BecameFollower)
		rn.CurrentTerm = peer.CurrentTerm
		rn.VotedFor = peer.Id
		rn.ElectionInProgress = false
		rn.SendSignal(raft.ElectionTimerStopped)
	}
}

func (em *ElectionManager) becomeAFollowerAccordingToPeer(rn *raft.RaftNode, v *pb.VoteResponse) {
	if rn.ElectionInProgress == true {
		rn.CurrentTerm = v.Term
		rn.CurrentRole = raft.FOLLOWER
		rn.VotedFor = 0
		rn.SendSignal(raft.ElectionTimerStopped)
		rn.ElectionInProgress = false
	}
}

func (em *ElectionManager) startElectionTimer(rn *raft.RaftNode) {
	em.electionTimedOut = make(chan bool)
	go func() {
		rn.SendSignal(raft.ElectionTimerStarted)
		timer := time.NewTimer(em.ElectionTimeoutDuration)
		for {
			select {
			case <-timer.C:
				if rn.ElectionInProgress {
					rn.SendSignal(raft.ElectionTimerStopped)
					timer.Stop()
					em.electionTimedOut <- true
					return
				} else {
					rn.SendSignal(raft.ElectionTimerStopped)
					timer.Stop()
					return
				}
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

func (em *ElectionManager) GetLeaderHeartChannel() chan raft.RaftNode {
	return em.leaderHeartbeatChannel
}

func (em *ElectionManager) GetResponseForVoteRequest(rn *raft.RaftNode, vr *pb.VoteRequest) *pb.VoteResponse {
	rn.SendSignal(raft.VoteRequestReceived)
	voteResponse := em.getVoteResponseForVoteRequest(rn, vr)

	if voteResponse.VoteGranted {
		rn.SendSignal(raft.VoteGranted)
	} else {
		rn.SendSignal(raft.VoteNotGranted)
	}
	return voteResponse
}
