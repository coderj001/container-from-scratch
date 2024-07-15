package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/onsi/gomega/gexec"
)

var pathToNsProcessCLI string

func TestNsProcess(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ns-process Suite")
}

var _ = BeforeSuite(func() {
	var err error
	pathToNsProcessCLI, err = gexec.Build("ns-process")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
