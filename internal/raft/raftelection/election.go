package raftelection

import (
	"math/rand"
	"time"

	"github.com/xaraphix/Sif/internal/raft"
)

var (
	ElctnMgr raft.RaftElection = &ElectionManager{
		ElectionTimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
		ElectionTimerOff:        true,
		VotesReceived: []raft.VoteResponse{},
	}
	leaderHeartbeatChannel chan raft.Peer
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

func (em *ElectionManager) GetReceivedVotes() []raft.VoteResponse {
	return em.VotesReceived
}

func (em *ElectionManager) GetLeaderHeartChannel() chan raft.Peer {
	if leaderHeartbeatChannel == nil {
		leaderHeartbeatChannel = make(chan raft.Peer)
	}

	return leaderHeartbeatChannel
}

func (el *ElectionManager) HasElectionTimerStarted() bool {
	return !el.ElectionTimerOff
}

func (el *ElectionManager) HasElectionTimerStopped() bool {
	return el.ElectionTimerOff
}

func (em *ElectionManager) StartElection(rn *raft.RaftNode) {
	em.initChannels()
	em.becomeACandidate(rn)
	em.askForVotes(rn)
	em.startElectionTimer(rn)
	em.concludeFromReceivedVotes(rn)
	restartElection := em.handleElection(rn)

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
			em.becomeAFollowerAccordingToPeer(rn, vr)
			return false
		case <-em.leader:
			em.ElectionTimerOff = true
			em.becomeALeader(rn)
			return false
		case l := <-em.leaderHeartbeatChannel:
			em.ElectionTimerOff = true
			em.becomeAFollowerAccordingToLeader(l)
			return false
		case <-em.follower:
			em.ElectionTimerOff = true
			em.becomeAFollower(rn)
			return false
		case <-em.electionTimedOut:
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

func (em *ElectionManager) concludeFromReceivedVotes(rn *raft.RaftNode) {

	em.followerAccordingToPeer = make(chan raft.VoteResponse)
	em.leader = make(chan bool)
	em.follower = make(chan bool)

	go func() {
		for {
			select {
			case <-em.conclusionDone:
				break
			default:
				votes := []raft.VoteResponse{}
				for vr := range em.votesResponse {
					em.VotesReceived = append(em.VotesReceived, vr)
					em.concludeFromReceivedVote(rn, votes, vr)
				}
			}
		}
	}()
}

func (em *ElectionManager) concludeFromReceivedVote(rn *raft.RaftNode, votes []raft.VoteResponse, vr raft.VoteResponse) {
	votes = append(votes, vr)
	valid := isElectionValid(rn, vr)
	if !valid && isPeerTermHigher(rn, vr) {
		em.followerAccordingToPeer <- vr
	} else if becomeLeader, becomeFollower := whatDoesTheMajorityWant(len(rn.Peers), votes, vr); true {
		if becomeLeader {
			em.leader <- true
		} else if becomeFollower {
			em.follower <- true
		}
	}

}

func isPeerTermHigher(rn *raft.RaftNode, vr raft.VoteResponse) bool {
	return rn.CurrentTerm < vr.Term
}

func isElectionValid(rn *raft.RaftNode, vr raft.VoteResponse) bool {
	if rn.CurrentRole == raft.CANDIDATE &&
		rn.CurrentTerm == vr.Term {
		return true
	} else {
		return false
	}
}

func whatDoesTheMajorityWant(numOfPeers int, votesReceived []raft.VoteResponse, vr raft.VoteResponse) (bool, bool) {

	majorityCount := numOfPeers / 2
	votesInFavor := 0
	leader, follower := false, false
	for _, vr := range votesReceived {
		if vr.VoteGranted {
			votesInFavor = votesInFavor + 1
		}
	}

	// if majority has voted against, game over
	if len(votesReceived)-votesInFavor > majorityCount {
		follower = true
	}

	// if majority has voted in favor, game won
	if votesInFavor >= majorityCount {
		leader = true
	}

	return leader, follower
}

func (em *ElectionManager) becomeAFollowerAccordingToLeader(leader raft.Peer) {

}

func (em *ElectionManager) becomeAFollowerAccordingToPeer(rn *raft.RaftNode, v raft.VoteResponse) {
	rn.CurrentTerm = v.Term
	rn.CurrentRole = raft.FOLLOWER
	rn.VotedFor = 0
	rn.IsHeartBeating = false
}

func (em *ElectionManager) restartElection() {
	
}

func (em *ElectionManager) startElectionTimer(rn *raft.RaftNode) {
	em.electionTimedOut = make(chan bool)
	go func() {
		em.ElectionTimerOff = false
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

func (em *ElectionManager) GetResponseForVoteRequest(rn *raft.RaftNode, vr raft.VoteRequest) raft.VoteResponse {
	return getVoteResponseForVoteRequest(rn, vr)
}
