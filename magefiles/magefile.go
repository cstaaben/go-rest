//go:build mage

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

func init() {
	os.Setenv("MAGEFILE_ENABLE_COLOR", "true")
}

var (
	Default = Build

	goCmd = mg.GoCmd()
)

// Build the application. If no source files have changed since the binary was built, it will not run the build command.
func Build() error {
	modified, err := target.Dir("./bin/go-rest", "./cmd", "./internal")
	if err != nil {
		log.Println("error checking file timestamps:", err)
	}

	if !modified {
		fmt.Println("no files changed, skipping build step")
		return nil
	}

	mg.SerialDeps(Tidy, Vendor, Lint)
	return sh.RunV(goCmd, "build", "-o", "bin/", "./cmd/go-rest")
}

// 'go run' with the example config
func Run() error {
	_, err := os.Stat("./logs")
	if os.IsNotExist(err) {
		if err = os.Mkdir("./logs", 0755); err != nil {
			return fmt.Errorf("creating logs dir: %w", err)
		}
	}

	mg.SerialDeps(Tidy, Vendor, Lint)
	return sh.RunV(goCmd, "run", "./cmd/go-rest", "--config", "./example/config.yaml")
}

// Run golangci-lint
func Lint() error {
	mg.SerialDeps(Tidy, Vendor)
	return sh.RunV("golangci-lint", "run", "--allow-parallel-runners")
}

// Tidy go modules
func Tidy() error {
	return sh.RunV(goCmd, "mod", "tidy")
}

// Download go modules
func Vendor() error {
	return sh.RunV(goCmd, "mod", "vendor")
}

// Kill any running binary with the filename from the build target.
func Kill() error {
	return sh.RunV("pkill", "--signal", "9", "--echo", "go-rest")
}
