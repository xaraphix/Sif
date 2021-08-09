package raftelection

import (
	"math/rand"
	"time"

	"github.com/xaraphix/Sif/internal/raft"
)


type ElectionManager struct {
	ElectionTimeoutDuration time.Duration
	ElectionTimerOff        bool

	VotesReceived []raft.VoteResponse

	electionTimerDone       chan bool
	votesResponse           chan raft.VoteResponse
	conclusionDone          chan bool
	electionTimedOut        chan bool
	followerAccordingToPeer chan raft.VoteResponse
	leader                  chan bool
	follower                chan bool
	leaderHeartbeatChannel  chan raft.Peer
}

func NewElectionManager() raft.RaftElection {
	return &ElectionManager{
		ElectionTimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
		ElectionTimerOff:        true,
		VotesReceived:           []raft.VoteResponse{},
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
		em.StartElection(rn)
	} else {
		rn.ElectionInProgress = false
		em.ElectionTimerOff = true
	}

}
func (em *ElectionManager) initChannels() {
	em.electionTimerDone = make(chan bool)
	em.votesResponse = make(chan raft.VoteResponse)
	em.conclusionDone = make(chan bool)
	em.electionTimedOut = make(chan bool)
	em.followerAccordingToPeer = make(chan raft.VoteResponse)
	em.leader = make(chan bool)
	em.follower = make(chan bool)
	em.leaderHeartbeatChannel = make(chan raft.Peer)
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
			em.ElectionTimerOff = true
			rn.SendSignal(raft.ElectionTimerStopped)
			em.becomeAFollowerAccordingToPeer(rn, vr)
			return false
		case <-em.leader:
			em.ElectionTimerOff = true
			rn.SendSignal(raft.ElectionTimerStopped)
			em.becomeALeader(rn)
			return false
		case l := <-em.leaderHeartbeatChannel:
			em.ElectionTimerOff = true
			rn.SendSignal(raft.ElectionTimerStopped)
			em.becomeAFollowerAccordingToLeader(l)
			return false
		case <-em.follower:
			em.ElectionTimerOff = true
			rn.SendSignal(raft.ElectionTimerStopped)
			em.becomeAFollower(rn)
			return false
		case <-em.electionTimedOut:
			rn.SendSignal(raft.ElectionTimerStopped)
			em.ElectionTimerOff = true
			return true
		}
	}
}

func (em *ElectionManager) becomeAFollower(rn *raft.RaftNode) {
	rn.CurrentRole = raft.FOLLOWER
	rn.ElectionInProgress = false
	rn.IsHeartBeating = false
}

func (em *ElectionManager) becomeALeader(rn *raft.RaftNode) {
	rn.CurrentRole = raft.LEADER
	em.replicateLogs(rn)
	rn.Heart.StartBeating(rn)
}

func (em *ElectionManager) replicateLogs(rn *raft.RaftNode) {
	go func(n *raft.RaftNode) {
		for _, peer := range rn.Peers {
			rn.LogMgr.ReplicateLog(rn, peer)
		}
	}(rn)
}

func (em *ElectionManager) becomeAFollowerAccordingToLeader(leader raft.Peer) {

}

func (em *ElectionManager) becomeAFollowerAccordingToPeer(rn *raft.RaftNode, v raft.VoteResponse) {
	rn.CurrentTerm = v.Term
	rn.CurrentRole = raft.FOLLOWER
	rn.VotedFor = 0
	rn.IsHeartBeating = false
}

func (em *ElectionManager) startElectionTimer(rn *raft.RaftNode) {
	em.electionTimedOut = make(chan bool)
	go func() {
		em.ElectionTimerOff = false
		rn.SendSignal(raft.ElectionTimerStarted)
		for {
			select {
			case <-em.electionTimerDone:
				break
			default:
				time.Sleep(em.ElectionTimeoutDuration)
				em.electionTimedOut <- true
			}
		}
	}()
}

func (em *ElectionManager) askForVotes(rn *raft.RaftNode) {
	voteRequest := em.GenerateVoteRequest(rn)
	for _, peer := range rn.Peers {
		go func(p raft.Peer) {
			voteResponse := rn.RPCAdapter.RequestVoteFromPeer(p, voteRequest)
			em.votesResponse <- voteResponse
		}(peer)
	}
}

func (em *ElectionManager) GetReceivedVotes() []raft.VoteResponse {
	return em.VotesReceived
}

func (em *ElectionManager) GetLeaderHeartChannel() chan raft.Peer {
	if em.leaderHeartbeatChannel == nil {
		em.leaderHeartbeatChannel = make(chan raft.Peer)
	}

	return em.leaderHeartbeatChannel
}

func (el *ElectionManager) HasElectionTimerStarted() bool {
	return !el.ElectionTimerOff
}

func (el *ElectionManager) HasElectionTimerStopped() bool {
	return el.ElectionTimerOff
}

func (em *ElectionManager) GetResponseForVoteRequest(rn *raft.RaftNode, vr raft.VoteRequest) raft.VoteResponse {
	return getVoteResponseForVoteRequest(rn, vr)
}
