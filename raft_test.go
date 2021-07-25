package raft

import (
	"reflect"
	"testing"
	"time"
)

func TestNodeShouldBecomeAFollowerOnInitialize(t *testing.T) {
	node, err := CreateRaftNode()

	if err != nil {
		t.Error("Node could not be initialized")
	}

	assertEqual(t, node.CurrentRole, FOLLOWER)
}

func TestNodeShouldBecomeCandidateOnLeaderHeartBeatTimeout(t *testing.T) {

	node, _ := CreateRaftNode()
	assertEqual(t, node.CurrentRole, FOLLOWER)

	node.LeaderLastHeartbeat = time.Now()

	time.Sleep(node.HeartbeatTimeout)

	assertEqual(t, node.CurrentRole, CANDIDATE)

}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
