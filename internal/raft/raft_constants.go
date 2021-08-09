package raft

const (
	FOLLOWER  = "follower"
	CANDIDATE = "candidate"
	LEADER    = "leader"

	ConfigLoaded = iota
	ConfigFound
	ConfigNotFound
	LockFileFound
	LockFileNotFound
	StartingFromCrash
	StartedFromCrash
	StartingClean
	StartedClean
	NodeInitialized
	LeaderHeartbeatMonitorStarted
	LeaderHeartbeatMonitorStopped
	ElectionStarted
	ElectionStopped
	ElectionTimerStarted
	ElectionTimerStopped
	ElectionRestarted
	RequestedVotes
	ReceivedVoteResponse
	BecameFollower
	BecameLeader
	VoteResponseSent
	LogReplicationRequestSent
)

