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
	LockFile            string
	RaftInstanceDirPath string = "/var/lib/sif/"
	RaftCfg             Config
	RaftVersion         string = "v0.1"
)

func NewConfig() Config {
	RaftCfg = Config{}
	RaftCfg.LoadConfig()

	return RaftCfg
}

type Config struct {
	RaftInstanceName    string      `yaml: "name"`
	RaftInstanceId      int32       `yaml: "id"`
	RaftPeers           []raft.Peer `yaml: peers`
	RaftInstanceDirPath string      `yaml: "sifdir"`
	RaftVersion         string      `yaml: "version"`

	BootedFromCrash bool
}

func (c Config) LoadConfig() {

	c = parseConfigFile()

	c.RaftPeers = getOrDefault(c.Peers, nil).([]raft.Peer)
	c.RaftInstanceDirPath = getOrDefault(c.RaftInstanceDirPath, RaftInstanceDirPath).(string)
	c.RaftInstanceId = getOrDefault(c.RaftInstanceId, 0).(int32)
	c.RaftInstanceName = getOrDefault(c.RaftInstanceName, getDefaultName()).(string)
	c.RaftVersion = RaftVersion
	c.BootedFromCrash = c.DidNodeCrash()
}

func (c Config) YamlFile() ([]byte, error) {
	filename, _ := filepath.Abs("./sifconfig.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	return yamlFile, err
}

func getDefaultName() string {
	return "Sif-" + strconv.Itoa(rand.Int())
}

func parseConfigFile() Config {
	cfg := Config{}

	yamlFile, err := cfg.YamlFile()
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func (c Config) DidNodeCrash() bool {

	filename, _ := filepath.Abs(c.RaftInstanceDirPath + ".siflock")
	_, err := ioutil.ReadFile(filename)

	if err != nil {
		return false
	}

	return true
}

func (c Config) InstanceDirPath() string {
	return c.RaftInstanceDirPath
}

func (c Config) InstanceId() int32 {
	return c.RaftInstanceId
}

func (c Config) InstanceName() string {
	return c.RaftInstanceName
}

func (c Config) Peers() []raft.Peer {
	return c.RaftPeers
}

func (c Config) Version() string {
	return c.RaftVersion
}
func getOrDefault(prop interface{}, defaultVal interface{}) interface{} {

	if prop == nil {
		return defaultVal
	}

	return prop
}
