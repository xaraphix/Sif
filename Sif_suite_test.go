package raft_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSif(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sif Suite")
}
