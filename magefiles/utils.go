package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/magefile/mage/sh"
)

// RunInPath runs a command in a specific path.
func runInPath(path string, cmd string, args ...string) (err error) {
	oldPath, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(path)
	if err != nil {
		return err
	}

	defer func() {
		err = os.Chdir(oldPath)
	}()

	return sh.Run(cmd, args...)
}

func outputInPath(path string, cmd string, args ...string) (_ string, err error) {
	oldPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	err = os.Chdir(path)
	if err != nil {
		return "", err
	}

	defer func() {
		err = os.Chdir(oldPath)
	}()

	return sh.Output(cmd, args...)
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func parseVersion(tag string) (*Version, error) {
	pattern := `v(?P<Major>0|(?:[1-9]\d*))(?:\.(?P<Minor>0|(?:[1-9]\d*))(?:\.(?P<Patch>0|(?:[1-9]\d*)))?(?:\-(?P<PreRelease>[0-9A-Z\.-]+))?(?:\+(?P<Meta>[0-9A-Z\.-]+))?)?`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}

	matches := re.FindStringSubmatch(tag)
	if matches == nil {
		return nil, fmt.Errorf("tag version '%s' is not valid, should be 'vX.Y.Z'", tag)
	}
	version := &Version{}
	for i, name := range re.SubexpNames() {
		switch name {
		case "Major":
			version.Major, err = strconv.Atoi(matches[i])
		case "Minor":
			version.Minor, err = strconv.Atoi(matches[i])
		case "Patch":
			version.Patch, err = strconv.Atoi(matches[i])
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse version %s: %w", name, err)
		}
	}
	return version, nil
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// EnsureGit ensures that git is installed, if not it panics.
func ensureGit() {
	if err := sh.Run("git", "--help"); err != nil {
		panic("git is not installed")
	}
}

// EnsureCargo ensures that cargo is installed, if not it panics.
func ensureCargo() {
	if err := sh.Run("cargo", "--help"); err != nil {
		panic("cargo is not installed")
	}
}

// EnsureCargoMake ensures that cargo-make is installed, if not, it tries to install it.
func ensureCargoMake() {
	ensureCargo()

	if err := sh.Run("cargo", "make", "--help"); err != nil {
		if err := sh.Run("cargo", "install", "cargo-make"); err != nil {
			panic(fmt.Sprintf("failed to install cargo-make: %v", err))
		}
	}
}

// EnsureQuicktype ensures that quicktype is installed, if not it panics.
func ensureQuicktype() {
	if err := sh.Run("quicktype", "--help"); err != nil {
		panic("quicktype is not installed")
	}
}

// EnsureYarn ensures that yarn is installed, if not it panics.
func ensureYarn() {
	if err := sh.Run("yarn", "--help"); err != nil {
		panic("yarn is not installed")
	}
}

// EnsureGoCodegen ensures that go-codegen is installed, if not it tries to install it.
func ensureGoCodegen() {
	if err := sh.Run("go-codegen", "--help"); err == nil {
		return
	}

	fmt.Printf("🔨 Installing go-codegen...\n")
	if err := sh.Run("go", "install", "github.com/srdtrk/go-codegen@0.2.5"); err != nil {
		panic(fmt.Sprintf("failed to install go-codegen: %v", err))
	}
}
