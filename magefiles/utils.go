//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"os"
	"path/filepath"
	"strings"
)

// RunInPath runs a command in a specific path.
func runInPath(path string, cmd string, args ...string) error {
	err := os.Chdir(path)
	if err != nil {
		return err
	}

	defer func() {
		dirs := filepath.SplitList(path)
		os.Chdir(strings.Repeat("../", len(dirs)))
	}()

	return sh.Run(cmd, args...)
}

// EnsureGit ensures that git is installed, if not it panics.
func ensureGit() {
	if err := sh.Run("command", "-v", "git"); err != nil {
		panic("git is not installed")
	}
}

// EnsureCargo ensures that cargo is installed, if not it panics.
func ensureCargo() {
	if err := sh.Run("command", "-v", "cargo"); err != nil {
		panic("cargo is not installed")
	}
}

// EnsureCargoMake ensures that cargo-make is installed, if not, it tries to install it.
func ensureCargoMake() {
	ensureCargo()

	if err := sh.Run("cargo", "make", "--help"); err == nil {
		return
	}

	sh.Run("cargo", "install", "--force", "cargo-make")
}

// EnsureQuicktype ensures that quicktype is installed, if not it panics.
func ensureQuicktype() {
	t, e := sh.Output("type", "quicktype")
	fmt.Printf("⚠️ : %s, error : %s\n", t, e)
	if err := sh.Run("type", "quicktype"); err != nil {
		panic("quicktype is not installed")
	}
}
