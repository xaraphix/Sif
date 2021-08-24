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
	VoteRequestReceived
	VoteGranted
	VoteNotGranted
	ReceivedVoteResponse
	BecameFollower
	BecameLeader
	BecameCandidate
	HeartbeatStarted
	HeartbeatStopped
	LogRequestSent
	ForwardedBroadcastReq
	DeliveredToApplication
	MsgAppendedToLogs
	LeaderHeartbeatReceived
)
