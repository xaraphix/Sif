package raftfile

import (
	"io/ioutil"
	"path/filepath"
)

type RaftFileMgr struct {
}

func NewFileManager() RaftFileMgr {
	return RaftFileMgr{}
}

func (fm RaftFileMgr) LoadFile(path string) ([]byte, error) {
	filename, _ := filepath.Abs(path)
	return ioutil.ReadFile(filename)
}

func (fm RaftFileMgr) SaveFile(path string) error {
	return nil
}
