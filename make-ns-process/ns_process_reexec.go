package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/pkg/reexec"
)

func init() {
	// Register the function to be called when the process is re-executed
	reexec.Register("nsInitialisation", nsInitialisation)
	// If reexec.Init() returns true, the process is a child re-execution
	if reexec.Init() {
		os.Exit(0)
	}
}

func nsInitialisation() {
	// Placeholder for namespace setup code
	fmt.Printf("\n>> namespace setup code goes here <<\n\n")
	nsRun()
}

func nsRun() {
	// Run a new shell in the namespace
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "PS1=-[ns-process]- # ")

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}

func main() {
	// Re-execute the current binary with nsInitialisation as the first argument
	cmd := reexec.Command("nsInitialisation")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error reexecuting command - %s\n", err)
		os.Exit(1)
	}
}
