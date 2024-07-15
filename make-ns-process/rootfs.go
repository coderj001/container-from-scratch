package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func pivotRoot(newroot string) error {
	// Ensure newroot is an absolute path
	newroot, err := filepath.Abs(newroot)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Check if newroot exists and is a directory
	if fi, err := os.Stat(newroot); err != nil || !fi.IsDir() {
		return fmt.Errorf("newroot %s does not exist or is not a directory", newroot)
	}

	// Create a temporary directory for the old root
	putold := filepath.Join(newroot, ".pivot_root")
	if err := os.MkdirAll(putold, 0700); err != nil {
		return fmt.Errorf("failed to create putold directory: %v", err)
	}

	// Bind mount newroot to itself
	if err := syscall.Mount(newroot, newroot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("failed to bind mount newroot: %v", err)
	}

	// Change to newroot directory
	if err := os.Chdir(newroot); err != nil {
		return fmt.Errorf("failed to change to newroot directory: %v", err)
	}

	// Perform pivot_root
	if err := syscall.PivotRoot(".", putold); err != nil {
		return fmt.Errorf("failed to pivot_root: %v", err)
	}

	// Change to root directory
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to change to root directory: %v", err)
	}

	// Unmount old root
	if err := syscall.Unmount(putold, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("failed to unmount old root: %v", err)
	}

	// Remove old root directory
	if err := os.Remove(putold); err != nil {
		return fmt.Errorf("failed to remove old root: %v", err)
	}

	return nil
}

func mountProc(newroot string) error {
	target := filepath.Join(newroot, "/proc")
	if err := os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("failed to create /proc directory: %v", err)
	}

	if err := syscall.Mount("proc", target, "proc", 0, ""); err != nil {
		return fmt.Errorf("failed to mount /proc: %v", err)
	}

	return nil
}

func exitIfRootfsNotFound(rootfsPath string) {
	if _, err := os.Stat(rootfsPath); os.IsNotExist(err) {
		usefulErrorMsg := fmt.Sprintf(`
"%s" does not exist.
Please create this directory and unpack a suitable root filesystem inside it.
An example rootfs, BusyBox, can be downloaded from:

https://raw.githubusercontent.com/teddyking/ns-process/4.0/assets/busybox.tar

And unpacked by:

mkdir -p %s
tar -C %s -xf busybox.tar
`, rootfsPath, rootfsPath, rootfsPath)

		fmt.Println(usefulErrorMsg)
		os.Exit(1)
	}
}
