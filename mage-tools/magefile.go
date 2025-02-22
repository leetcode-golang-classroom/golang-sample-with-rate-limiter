//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

// clean the build binary
func Clean() error {
	return sh.Rm("bin")
}

// install - update the dependency
func Install() error {
	return sh.Run("go", "mod", "download")
}

// build-server - build server
func BuildServer() error {
	// build server
	return sh.Run("go", "build", "-o", "./bin/server", "./cmd/server/main.go")
}

// build-client - build client
func BuildClient() error {
	// build server
	return sh.Run("go", "build", "-o", "./bin/client", "./cmd/client/main.go")
}

// Creates the binary in the current directory.
func Build() error {
	mg.Deps(Clean)
	mg.Deps(Install)
	mg.Deps(BuildClient)
	return BuildServer()
}

// start server
func StartServer() error {
	mg.Deps(BuildServer)
	return sh.RunV("./bin/server")
}

// start client
func StartClient() error {
	mg.Deps(BuildClient)
	return sh.RunV("./bin/client")
}

// run the test
func Test() error {
	err := sh.RunV("go", "test", "-v", "./...")
	if err != nil {
		return err
	}
	return nil
}
