//go:build (darwin && cgo) || linux || mage
// +build darwin,cgo linux mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

// Ts build typescript schema for the given contract schema.
func (Build) Ts(schema string) error {
	fmt.Printf("⚙️ Generate typescript types for %s\n", schema)

	ensureYarn()
	ensureQuicktype()

	name := strings.TrimPrefix(schema, "axone-")
	dest := filepath.Join(TS_DIR, fmt.Sprintf("%s-schema", name))
	if err := os.MkdirAll(filepath.Join(dest, "gen-ts"), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err := sh.Run("bash", "-c",
		fmt.Sprintf("quicktype -s schema %s -o %s --prefer-types --prefer-unions",
			filepath.Join(SCHEMA_DIR, schema, "*.json"),
			filepath.Join(dest, "gen-ts", "schema.ts")))
	if err != nil {
		return fmt.Errorf("failed to generate typescript types: %w", err)
	}

	fmt.Println("🔨 Building typescript")

	err = sh.Run("yarn", "--cwd", dest)
	if err != nil {
		return err
	}

	return sh.Run("yarn", "--cwd", dest, "build")
}

// Go build go schema for the given contract schema.
func (Build) Go(schema string) error {
	fmt.Printf("⚙️ Generate go types for %s\n", schema)

	ensureGoCodegen()

	name := strings.TrimPrefix(schema, "axone-")
	dest := filepath.Join(GO_DIR, fmt.Sprintf("%s-schema", name))
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err := sh.Run("go-codegen", "generate",
		"messages",
		filepath.Join(SCHEMA_DIR, fmt.Sprintf("%s.json", schema)),
		"-o", filepath.Join(dest, "schema.go"),
		"--package-name", "schema")
	if err != nil {
		return fmt.Errorf("failed to generate go types: %w", err)
	}

	fmt.Println("🔨 Building go")
	return runInPath(dest, "go", "build")
}
