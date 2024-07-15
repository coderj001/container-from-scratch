package main_test

import (
	"bufio"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

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

		fmt.Println("Setting up command")
		command = exec.Command("sudo", pathToNsProcessCLI)
		stdin, err = command.StdinPipe()
		Expect(err).NotTo(HaveOccurred())
		stdout = gbytes.NewBuffer()
		cmdToRunInNamespacedShell = "true"
	})

	JustBeforeEach(func() {
		var err error

		fmt.Println("Writing command to stdin")
		stdinWriter := bufio.NewWriter(stdin)
		stdinWriter.WriteString(cmdToRunInNamespacedShell)
		stdinWriter.Flush()
		Expect(stdin.Close()).To(Succeed())

		fmt.Println("Starting session")
		session, err = gexec.Start(command, stdout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		fmt.Println("Closing stdout")
		Expect(stdout.Close()).To(Succeed())
	})

	It("exits with a 0 exit code", func() {
		Eventually(session).Should(gexec.Exit(0))
	})

	Describe("namespace initialization", func() {
		BeforeEach(func() {
			cmdToRunInNamespacedShell = "cat /proc/self/mounts"
		})

		It("sets up namespaces correctly", func() {
			Eventually(stdout).Should(gbytes.Say("/"))
		})

		It("runs commands in the new namespace", func() {
			cmdToRunInNamespacedShell = "echo 'Hello from ns-process'"
			Eventually(stdout).Should(gbytes.Say("Hello from ns-process"))
		})
	})

	Describe("user namespace configuration", func() {
		BeforeEach(func() {
			cmdToRunInNamespacedShell = "id"
		})

		It("applies a UID mapping", func() {
			Eventually(stdout).Should(gbytes.Say(`uid=0\(root\)`))
		})

		It("applies a GID mapping", func() {
			Eventually(stdout).Should(gbytes.Say(`gid=0\(root\)`))
		})
	})
})

func inode(pid, namespaceType string) string {
	namespace, err := os.Readlink(fmt.Sprintf("/proc/%s/ns/%s", pid, namespaceType))
	Expect(err).NotTo(HaveOccurred())

	requiredFormat := regexp.MustCompile(`^\w+:\[\d+\]$`)
	Expect(requiredFormat.MatchString(namespace)).To(BeTrue())

	namespace = strings.Split(namespace, ":")[1]
	namespace = namespace[1:]
	namespace = namespace[:len(namespace)-1]

	return namespace
}
