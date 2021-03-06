package testbed_setup

import (
	"time"

	"github.com/xaraphix/Sif/internal/raft"
	"github.com/xaraphix/Sif/internal/raft/protos"
	"github.com/xaraphix/Sif/internal/raft/raftadapter"
	"github.com/xaraphix/Sif/internal/raft/raftconfig"
	"github.com/xaraphix/Sif/internal/raft/raftelection"
	"github.com/xaraphix/Sif/internal/raft/raftfile"
	"github.com/xaraphix/Sif/internal/raft/raftlog"
)

func Setup5FollowerNodes() (*raft.RaftNode, *raft.RaftNode, *raft.RaftNode, *raft.RaftNode, *raft.RaftNode) {

	currDirPath := GetCurrentDirPath()

	deps1 := GetNewNodeDeps()
	deps1.ConfigManager.SetConfigFilePath(currDirPath + "data/node1_config.yml")
	node1 := raft.NewRaftNode(deps1)

	deps2 := GetNewNodeDeps()
	deps2.ConfigManager.SetConfigFilePath(currDirPath + "data/node2_config.yml")
	node2 := raft.NewRaftNode(deps2)

	deps3 := GetNewNodeDeps()
	deps3.ConfigManager.SetConfigFilePath(currDirPath + "data/node3_config.yml")
	node3 := raft.NewRaftNode(deps3)

	deps4 := GetNewNodeDeps()
	deps4.ConfigManager.SetConfigFilePath(currDirPath + "data/node4_config.yml")
	node4 := raft.NewRaftNode(deps4)

	deps5 := GetNewNodeDeps()
	deps5.ConfigManager.SetConfigFilePath(currDirPath + "data/node5_config.yml")
	node5 := raft.NewRaftNode(deps5)
	return node1, node2, node3, node4, node5
}

func GetNewNodeDeps() raft.RaftDeps {
	options := raft.RaftOptions{
		StartLeaderHeartbeatMonitorAfterInitializing: false,
	}

	return raft.RaftDeps{
		FileManager:      raftfile.NewFileManager(),
		ConfigManager:    raftconfig.NewConfig(),
		ElectionManager:  raftelection.NewElectionManager(),
		HeartbeatMonitor: raft.NewLeaderHeartbeatMonitor(true),
		RPCAdapter:       raftadapter.NewRaftNodeAdapter(),
		LogManager:       raftlog.NewLogManager(),
		Heart:            raftelection.NewLeaderHeart(),
		Options:          options,
	}

}

func ProceedWhenRPCAdapterStarted(nodes []**raft.RaftNode) {
	for {
		peer1 := (*nodes[0]).RPCAdapter.GetRaftInfo((*nodes[0]).Peers[0], &protos.RaftInfoRequest{})
		peer2 := (*nodes[0]).RPCAdapter.GetRaftInfo((*nodes[0]).Peers[1], &protos.RaftInfoRequest{})
		peer3 := (*nodes[0]).RPCAdapter.GetRaftInfo((*nodes[0]).Peers[2], &protos.RaftInfoRequest{})
		peer4 := (*nodes[0]).RPCAdapter.GetRaftInfo((*nodes[0]).Peers[3], &protos.RaftInfoRequest{})
		if peer1 != nil && peer2 != nil && peer3 != nil && peer4 != nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func ProceedWhenLeaderAccepted(nodes []**raft.RaftNode, leaderId string) {
  for i, n := range nodes {
    if i == 0 {
      continue
    }
    CheckIfEventTriggered((*n), raft.LeaderAccepted, raft.RaftEventDetails{CurrentLeader: leaderId})
  }
}

func ProceedLogAckReceived(nodes []**raft.RaftNode, leaderId string) {
  for i, n := range nodes {
    if i == 0 {
      continue
    }
    CheckIfEventTriggered((*nodes[0]), raft.AckLengthUpdated, raft.RaftEventDetails{Peer: (*n).Id, AckLength: int32(2)})
  }
}

func ProceedWhenLeaderCommitsLogs(node *raft.RaftNode) {
	CheckIfEventTriggered(node, raft.DeliveredToApplication, raft.RaftEventDetails{})
}

func ProceedWhenFollowersCommitLogs(nodes []**raft.RaftNode, cl int32) {
  CheckIfEventTriggered((*nodes[0]), raft.CommitLengthUpdated, raft.RaftEventDetails{Id: (*nodes[0]).Id, CommitLength: cl})
  for i, n := range nodes {
    if i == 0 {
      continue
    }
    CheckIfEventTriggered((*n), raft.CommitLengthUpdated, raft.RaftEventDetails{Id: (*n).Id, CommitLength: cl})
  }
}

func DestructAllNodes(nodes []**raft.RaftNode) {
	for i := range nodes {
		(*nodes[i]).Close()
		*nodes[i] = nil
	}
}
