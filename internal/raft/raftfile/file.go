package raftfile

import (
	"io/ioutil"
	"path/filepath"
)

var fileMgr *RaftFileMgr

func init() {
	fileMgr = &RaftFileMgr{}
}

type RaftFileMgr struct {
}

func NewRaftFileMfg() *RaftFileMgr {
	return fileMgr
}

func (fm *RaftFileMgr) LoadFile(path string) ([]byte, error) {
	filename, _ := filepath.Abs(path)
	return ioutil.ReadFile(filename)
}

func (fm *RaftFileMgr) SaveFile(path string) error {
	return nil
}
