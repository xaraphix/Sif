package raftconfig

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/xaraphix/Sif/internal/raft"
	"gopkg.in/yaml.v2"
)

var (
	LogsFile            string
	LockFile            string = "./.siflock"
	RaftInstanceDirPath string = "/var/lib/sif/"
	RaftCfg             Config
	RaftVersion         string = "v0.1"
)

type Config struct {
	RaftInstanceName    string      `yaml:"name"`
	RaftInstanceId      int32       `yaml:"id"`
	RaftPeers           []raft.Peer `yaml:"peers"`
	RaftInstanceDirPath string      `yaml:"sifdir"`
	RaftVersion         string      `yaml:"version"`

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
	c.RaftVersion = RaftVersion
	c.BootedFromCrash = c.DidNodeCrash(rn)
}

func (c *Config) LoadConfig(rn *raft.RaftNode) {
	parseConfig(c, rn)
}

func yamlFile() ([]byte, error) {
	filename, _ := filepath.Abs("./sifconfig.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	return yamlFile, err
}

func getDefaultName() string {
	return "Sif-" + strconv.Itoa(rand.Int())
}

func (c *Config) DidNodeCrash(rn *raft.RaftNode) bool {
	_, err := rn.FileMgr.LoadFile(LockFile)
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

func (c *Config) Version() string {
	return c.RaftVersion
}
func getOrDefault(prop interface{}, defaultVal interface{}) interface{} {

	if prop == nil {
		return defaultVal
	}

	return prop
}
