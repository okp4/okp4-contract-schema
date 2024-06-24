package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	CONTRACTS_REPOSITORY = "https://github.com/axone-protocol/contracts.git"
	CONTRACTS_TMP_DIR    = "tmp"
)

type Schema mg.Namespace

// Download download contracts schemas at a given ref.
// ref can be a branch, tag or commit hash.
func (Schema) Download(ref string) error {
	mg.Deps(Schema.Clean)

	fmt.Println("📥 Downloading contracts schemas")

	ensureGit()

	if err := sh.Run("git",
		"clone",
		"--depth", "1",
		"--branch", ref,
		CONTRACTS_REPOSITORY,
		CONTRACTS_TMP_DIR); err != nil {
		return err
	}

	return nil
}

// Generate build and generate contracts json schemas at the given ref.
// It first downloads the contracts repository at the given ref., then
// generate the json schema using cargo make and finally move the generated schema to the schema directory.
// ref can be a branch, tag or commit hash.
func (s Schema) Generate(ref string) error {
	mg.Deps(mg.F(Schema.Download, ref))

	defer func() { _ = s.Clean() }()

	fmt.Println("🔨 Generating contracts json schema")

	ensureCargoMake()

	if err := runInPath(CONTRACTS_TMP_DIR, "cargo", "make", "schema"); err != nil {
		return fmt.Errorf("failed to generate contracts json schema: %w", err)
	}

	if err := sh.Rm(SCHEMA_DIR); err != nil {
		return fmt.Errorf("failed to remove schema directory: %w", err)
	}

	fmt.Println("🔨 Moving generated json schema")
	if err := sh.Run("bash", "-c",
		fmt.Sprintf("rsync -rmv --include='*/' --include='*/schema/raw/*.json'  --include='*/schema/*.json' --exclude='*' %s/contracts/ %s/", CONTRACTS_TMP_DIR, SCHEMA_DIR)); err != nil {
		return fmt.Errorf("failed to move generated json schema: %w", err)
	}

	schemas, err := sh.Output("find", SCHEMA_DIR, "-type", "f", "-name", "*.json")
	if err != nil {
		return fmt.Errorf("failed to find all json schema in '%s': %w", SCHEMA_DIR, err)
	}

	for _, schema := range strings.Split(schemas, "\n") {
		dest := filepath.Join(schema, "../../../", filepath.Base(schema))
		if err := os.Rename(schema, dest); err != nil {
			return fmt.Errorf("failed to move schema '%s' to '%s': %w", schema, dest, err)
		}
	}
	fmt.Println("✨ Contracts json schema generated")
	return nil
}

// Readme generate contracts readme on all target.
func (s Schema) Readme(ref string) error {
	mg.Deps(mg.F(Schema.Download, ref))
	defer func() { _ = s.Clean() }()

	fmt.Println("📄 Generating contracts readme")

	contractsDir, err := os.ReadDir(SCHEMA_DIR)
	if err != nil {
		return fmt.Errorf("failed to read contracts directory '%s': %w", SCHEMA_DIR, err)
	}

	for _, contract := range contractsDir {
		if !contract.IsDir() {
			continue
		}

		err := generateReadmes(contract.Name())
		if err != nil {
			return fmt.Errorf("failed to generate readme for contract '%s': %w", contract.Name(), err)
		}
	}
	return nil
}

// tag returns the git tag for the current branch or "" if none.
func tag() string {
	s, _ := sh.Output("git", "describe", "--tags")
	return s
}

// Clean remove temporary files.
func (Schema) Clean() error {
	fmt.Println("🧹 Cleaning schema temporary files")
	return sh.Rm(CONTRACTS_TMP_DIR)
}
