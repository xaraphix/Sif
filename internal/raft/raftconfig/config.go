package raftconfig

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/xaraphix/Sif/internal/raft"
	"gopkg.in/yaml.v2"
)

var (
	PersistentStateFile string = "raft_state.json"
	LockFile            string = ".siflock"
	RaftInstanceDirPath string = "/var/lib/sif/"
	RaftCfg             Config
	RaftVersion         string = "v0.1"
)

type RaftPersistentState struct {
	RaftLogs         *[]raft.Log `json:"logs"`
	RaftCurrentTerm  int32       `json:"currentTerm"`
	RaftVotedFor     int32       `json:"votedFor"`
	RaftCommitLength int32       `json:"commitLength"`
}

type Config struct {
	RaftInstanceName                    string      `yaml:"name"`
	RaftInstanceId                      int32       `yaml:"id"`
	RaftPeers                           []raft.Peer `yaml:"peers"`
	RaftInstanceDirPath                 string      `yaml:"sifdir"`
	RaftVersion                         string      `yaml:"version"`
	RaftInstancePersistentStateFilePath string

	RaftPersistentState

	BootedFromCrash bool
}

func NewConfig() *Config {
	RaftCfg = Config{}
	return &RaftCfg
}

func parseConfig(c *Config, rn *raft.RaftNode) {
	cfgFile, err := rn.FileMgr.LoadFile("./sifconfig.yml")
	cfg := &Config{}
	err = yaml.Unmarshal(cfgFile, cfg)
	if err != nil {
		panic(err)
	}

	c.RaftPeers = getOrDefault(cfg.Peers(), nil).([]raft.Peer)
	c.RaftInstanceDirPath = getOrDefault(cfg.InstanceDirPath(), RaftInstanceDirPath).(string)
	c.RaftInstanceId = getOrDefault(cfg.InstanceId(), 0).(int32)
	c.RaftInstanceName = getOrDefault(cfg.InstanceName(), getDefaultName()).(string)
	c.RaftInstancePersistentStateFilePath = getOrDefault(cfg.LogFilePath(), getDefaultName()).(string)
	c.RaftInstancePersistentStateFilePath = getOrDefault(cfg.InstanceDirPath()+PersistentStateFile, getDefaultName()).(string)
	c.RaftVersion = RaftVersion
	c.BootedFromCrash = c.DidNodeCrash(rn)
	persistentState := loadRaftPersistentState(c.RaftInstancePersistentStateFilePath, rn)
	setPersistentState(rn, c, persistentState)
}

func setPersistentState(rn *raft.RaftNode, c *Config, s *RaftPersistentState) {
	if c.BootedFromCrash {
		c.RaftCurrentTerm = s.RaftCurrentTerm
		c.RaftCommitLength = s.RaftCommitLength
		c.RaftVotedFor = s.RaftVotedFor
		c.RaftLogs = s.RaftLogs
	}
}

func (c *Config) LoadConfig(rn *raft.RaftNode) {
	parseConfig(c, rn)
}

func loadConfigFile() ([]byte, error) {
	filename, _ := filepath.Abs("./sifconfig.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	return yamlFile, err
}

func getDefaultName() string {
	return "Sif-" + strconv.Itoa(rand.Int())
}

func loadRaftPersistentState(filepath string, rn *raft.RaftNode) *RaftPersistentState {
	stateFile, err := rn.FileMgr.LoadFile(filepath)
	state := &RaftPersistentState{}

	if err != nil {
		//TODO do something
	}

	err = json.Unmarshal(stateFile, state)

	if err != nil {
		//TODO do something
	}

	return state
}

func (c *Config) DidNodeCrash(rn *raft.RaftNode) bool {
	_, err := rn.FileMgr.LoadFile(c.RaftInstanceDirPath + LockFile)
	if err != nil {
		return false
	}

	return true
}

func (c *Config) InstanceDirPath() string {
	return c.RaftInstanceDirPath
}

func (c *Config) InstanceId() int32 {
	return c.RaftInstanceId
}

func (c *Config) InstanceName() string {
	return c.RaftInstanceName
}

func (c *Config) Peers() []raft.Peer {
	return c.RaftPeers
}

func (c *Config) LogFilePath() string {
	return c.RaftInstancePersistentStateFilePath
}

func (c *Config) CurrentTerm() int32 {
	return c.RaftCurrentTerm
}

func (c *Config) Logs() *[]raft.Log {
	return c.RaftLogs
}

func (c *Config) VotedFor() int32 {
	return c.RaftVotedFor
}

func (c *Config) CommitLength() int32 {
	return c.RaftCommitLength
}

func (c *Config) Version() string {
	return c.RaftVersion
}

func getOrDefault(prop interface{}, defaultVal interface{}) interface{} {

	if prop == nil {
		return defaultVal
	}

	return prop
}
