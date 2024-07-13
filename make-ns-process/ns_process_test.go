package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bufio"
	"io"
	"os/exec"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var pathToNsProcessCLI string

var _ = BeforeSuite(func() {
	var err error
	pathToNsProcessCLI, err = gexec.Build("ns-process")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("The ns-process CLI", func() {
	var (
		command                   *exec.Cmd
		session                   *gexec.Session
		stdin                     io.WriteCloser
		cmdToRunInNamespacedShell string
		stdout                    *gbytes.Buffer
	)

	BeforeEach(func() {
		var err error

		command = exec.Command(pathToNsProcessCLI)
		stdin, err = command.StdinPipe()
		Expect(err).NotTo(HaveOccurred())
		stdout = gbytes.NewBuffer()
		cmdToRunInNamespacedShell = "true"
	})

	JustBeforeEach(func() {
		var err error

		stdinWriter := bufio.NewWriter(stdin)
		stdinWriter.WriteString(cmdToRunInNamespacedShell)
		stdinWriter.Flush()
		Expect(stdin.Close()).To(Succeed())

		session, err = gexec.Start(command, stdout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(stdout.Close()).To(Succeed())
	})

	It("exits with a 0 exit code", func() {
		Eventually(session).Should(gexec.Exit(0))
	})
})
