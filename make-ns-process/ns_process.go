package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)

func init() {
	reexec.Register("nsInitialisation", nsInitialisation)
	if reexec.Init() {
		os.Exit(0)
	}
}

func nsInitialisation() {
	newRootPath := os.Args[1]

	if err := mountProc(newRootPath); err != nil {
		fmt.Printf("Error mounting /proc - %s\n", err)
		os.Exit(1)
	}

	if err := pivotRoot(newRootPath); err != nil {
		fmt.Printf("Error running pivot_root - %s\n", err)
		os.Exit(1)
	}
	nRun()
}

func nRun() {
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "PS1=- [ns-process] - #")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running /bin/sh command - %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("This program needs to be run as root.")
		os.Exit(1)
	}

	rootfsPath := "/tmp/ns-process/rootfs"

	cmd := reexec.Command("nsInitialisation", rootfsPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error reexecuting command - %s\n", err)
		os.Exit(1)
	}

}
