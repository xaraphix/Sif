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
	askForVotesDone         chan bool
	electionTimedOut        chan bool
	followerAccordingToPeer chan raft.VoteResponse
	leader                  chan bool
	follower                chan bool
	leaderHeartbeatChannel  chan raft.RaftNode
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
		rn.SendSignal(raft.ElectionRestarted)
	}

}
func (em *ElectionManager) initChannels() {
	em.electionTimerDone = make(chan bool)
	em.votesResponse = make(chan raft.VoteResponse)
	em.conclusionDone = make(chan bool)
	em.askForVotesDone = make(chan bool)
	em.electionTimedOut = make(chan bool)
	em.followerAccordingToPeer = make(chan raft.VoteResponse)
	em.leader = make(chan bool)
	em.follower = make(chan bool)
	em.leaderHeartbeatChannel = make(chan raft.RaftNode)
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
			rn.SendSignal(raft.ElectionTimerStopped)
			em.becomeAFollowerAccordingToPeer(rn, vr)
			return false
		case <-em.leader:
			rn.SendSignal(raft.ElectionTimerStopped)
			em.becomeALeader(rn)
			return false
		case l := <-em.leaderHeartbeatChannel:
			rn.SendSignal(raft.ElectionTimerStopped)
			em.becomeAFollowerAccordingToLeader(rn, l)
			return false
		case <-em.follower:
			rn.SendSignal(raft.ElectionTimerStopped)
			em.becomeAFollower(rn)
			return false
		case <-em.electionTimedOut:
			rn.SendSignal(raft.ElectionTimerStopped)
			return true
		}
	}
}

func (em *ElectionManager) becomeAFollower(rn *raft.RaftNode) {
	em.conclusionDone <- true
	em.askForVotesDone <- true
	em.electionTimerDone <- true
	rn.CurrentRole = raft.FOLLOWER
	rn.ElectionInProgress = false
	rn.IsHeartBeating = false
}

func (em *ElectionManager) becomeALeader(rn *raft.RaftNode) {
	em.conclusionDone <- true
	em.askForVotesDone <- true
	em.electionTimerDone <- true
	rn.CurrentRole = raft.LEADER
	rn.SendSignal(raft.BecameLeader)
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

func (em *ElectionManager) becomeAFollowerAccordingToLeader(rn *raft.RaftNode, leader raft.RaftNode) {
	em.conclusionDone <- true
	em.askForVotesDone <- true
	em.electionTimerDone <- true
	rn.CurrentRole = raft.FOLLOWER
	rn.CurrentTerm = leader.CurrentTerm
}

func (em *ElectionManager) becomeAFollowerAccordingToPeer(rn *raft.RaftNode, v raft.VoteResponse) {
	em.conclusionDone <- true
	em.askForVotesDone <- true
	em.electionTimerDone <- true

	rn.CurrentTerm = v.Term
	rn.CurrentRole = raft.FOLLOWER
	rn.VotedFor = 0
	rn.IsHeartBeating = false
}

func (em *ElectionManager) startElectionTimer(rn *raft.RaftNode) {
	em.electionTimedOut = make(chan bool)
	go func() {
		rn.SendSignal(raft.ElectionTimerStarted)
		timer := time.NewTimer(em.ElectionTimeoutDuration)
		for {
			select {
			case <-em.electionTimerDone:
				timer.Stop()
				break
			case <- timer.C:
				em.electionTimedOut <- true
			}
		}
	}()
}

func (em *ElectionManager) askForVotes(rn *raft.RaftNode) {
	voteRequest := em.GenerateVoteRequest(rn)
	for _, peer := range rn.Peers {
		go func(p raft.Peer) {

			for {
				select {
				case <-em.askForVotesDone:
					break
				default:

					voteResponse := rn.RPCAdapter.RequestVoteFromPeer(p, voteRequest)
					em.votesResponse <- voteResponse
					break
				}
			}
		}(peer)
	}
}

func (em *ElectionManager) GetReceivedVotes() []raft.VoteResponse {
	return em.VotesReceived
}

func (em *ElectionManager) GetLeaderHeartChannel() chan raft.RaftNode {
	if em.leaderHeartbeatChannel == nil {
		em.leaderHeartbeatChannel = make(chan raft.RaftNode)
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
	rn.SendSignal(raft.VoteRequestReceived)
	voteResponse := getVoteResponseForVoteRequest(rn, vr)

	if voteResponse.VoteGranted {
		rn.SendSignal(raft.VoteGranted)
	} else {
		rn.SendSignal(raft.VoteNotGranted)
	}
	return voteResponse
}
