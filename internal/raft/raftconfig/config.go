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

//go:generate mockgen -destination=../mocks/mock_configfile.go -package=mocks . File
type File interface {
	Load(filepath string) ([]byte, error)
}

type Config struct {
	RaftInstanceName    string      `yaml:"name"`
	RaftInstanceId      int32       `yaml:"id"`
	RaftPeers           []raft.Peer `yaml:"peers"`
	RaftInstanceDirPath string      `yaml:"sifdir"`
	RaftVersion         string      `yaml:"version"`

	BootedFromCrash bool

	configFile File
}

type ConfigFileIO struct {
	
}

func (cf *ConfigFileIO) Load(path string) ([]byte, error) {
	filename, _ := filepath.Abs(path)
	return ioutil.ReadFile(filename)
}

func NewConfig(cf File) raft.RaftConfig {
	if cf == nil {
		cf = &ConfigFileIO{}
	}

	RaftCfg = Config{
		configFile: cf,
	}
	return &RaftCfg
}

func (c *Config) LoadConfigFromFile() raft.RaftConfig {
	cfg := &Config{}

	yamlFile, err := c.configFile.Load("./sifconfig.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func (c *Config) InitializeConfig() {
	cfg := c.LoadConfigFromFile()

	c.RaftPeers = getOrDefault(cfg.Peers(), nil).([]raft.Peer)
	c.RaftInstanceDirPath = getOrDefault(cfg.InstanceDirPath(), RaftInstanceDirPath).(string)
	c.RaftInstanceId = getOrDefault(cfg.InstanceId(), 0).(int32)
	c.RaftInstanceName = getOrDefault(cfg.InstanceName(), getDefaultName()).(string)
	c.RaftVersion = RaftVersion
	c.BootedFromCrash = c.DidNodeCrash()
}

func yamlFile() ([]byte, error) {
	filename, _ := filepath.Abs("./sifconfig.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	return yamlFile, err
}

func getDefaultName() string {
	return "Sif-" + strconv.Itoa(rand.Int())
}

func (c *Config) DidNodeCrash() bool {

  _, err := c.configFile.Load("./.siflock")
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
