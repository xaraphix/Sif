package raftconfig

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/xaraphix/Sif/internal/raft"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v2"
)

var (
	PersistentStateFile string = "raft_state.json"
	LockFile            string = ".siflock"
	RaftInstanceDirPath string = "/var/lib/sif/"
	RaftCfg             Config
	RaftVersion         string = "v0.1"
)

type Config struct {
	RaftInstanceName                    string      `yaml:"name"`
	RaftInstanceId                      string       `yaml:"id"`
	RaftPeers                           []raft.Peer `yaml:"peers"`
	RaftInstanceDirPath                 string      `yaml:"sifdir"`
	RaftVersion                         string      `yaml:"version"`
	RaftInstancePersistentStateFilePath string      `yaml:"persistentStateFile"`
	RaftHost                            string      `yaml:"host"`
	RaftPort                            string      `yaml:"port"`
	ConfigFilePath                      string
	pb.RaftPersistentState

	BootedFromCrash bool
}

func NewConfig() *Config {
	return &Config{}
}

func parseConfig(c *Config, rn *raft.RaftNode) {

	path := getOrDefault(c.ConfigFilePath, "./sifconfig.yml").(string)
	cfgFile, err := rn.FileMgr.LoadFile(path)
	cfg := &Config{}
	err = yaml.Unmarshal(cfgFile, cfg)
	if err != nil {
		panic(err)
	}

	c.RaftHost = getOrDefault(cfg.Host(), nil).(string)
	c.RaftPort = getOrDefault(cfg.Port(), nil).(string)
	c.RaftPeers = getOrDefault(cfg.Peers(), nil).([]raft.Peer)
	c.RaftInstanceDirPath = getOrDefault(cfg.InstanceDirPath(), RaftInstanceDirPath).(string)
	c.RaftInstanceId = getOrDefault(cfg.InstanceId(), 0).(string)
	c.RaftInstanceName = getOrDefault(cfg.InstanceName(), getDefaultName()).(string)
	c.RaftInstancePersistentStateFilePath = getOrDefault(cfg.InstanceDirPath()+PersistentStateFile, getDefaultName()).(string)
	c.RaftVersion = RaftVersion
	c.BootedFromCrash = c.DidNodeCrash(rn)
	persistentState := loadRaftPersistentState(c.RaftInstancePersistentStateFilePath, rn)
	setPersistentState(rn, c, persistentState)
}

func setPersistentState(rn *raft.RaftNode, c *Config, s *pb.RaftPersistentState) {
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

func loadRaftPersistentState(filepath string, rn *raft.RaftNode) *pb.RaftPersistentState {
	stateFile, err := rn.FileMgr.LoadFile(filepath)
	state := &pb.RaftPersistentState{}

	if err != nil {
		//TODO do something
	}

	err = protojson.Unmarshal(stateFile, state)

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

func (c *Config) InstanceId() string {
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

func (c *Config) Logs() []*pb.Log {
	return c.RaftLogs
}

func (c *Config) VotedFor() string {
	return c.RaftVotedFor
}

func (c *Config) CommitLength() int32 {
	return c.RaftCommitLength
}

func (c *Config) Version() string {
	return c.RaftVersion
}

func (c *Config) Host() string {
	return c.RaftHost
}

func (c *Config) Port() string {
	return c.RaftPort
}

func (c *Config) SetConfigFilePath(path string) {
	c.ConfigFilePath = path
}

func getOrDefault(prop interface{}, defaultVal interface{}) interface{} {

	if prop == nil || prop == "" {
		return defaultVal
	}

	return prop
}
