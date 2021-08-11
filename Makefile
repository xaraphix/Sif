.PHONY: protos

protos:
	protoc --go_out=. --go_opt=paths=source_relative  internal/raft/raftadapter/protos/adapter.proto

