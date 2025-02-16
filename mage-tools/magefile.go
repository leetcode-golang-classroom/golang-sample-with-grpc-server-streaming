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

// setup redis
func SetupRedis() error {
	return sh.RunV("docker", "compose", "up", "-d", "redis")
}

// build-server - build grpc server
func BuildServer() error {
	// build grpc server
	return sh.Run("go", "build", "-o", "./bin/grpc-server", "./cmd/server/main.go")
}

// build-client - build grpc client
func BuildClient() error {
	// build grpc client
	return sh.Run("go", "build", "-o", "./bin/grpc-client", "./cmd/client/main.go")
}

// build-publisher - build redis publisher
func BuildPublisher() error {
	// build publisher
	return sh.Run("go", "build", "-o", "./bin/publisher", "./cmd/publisher/main.go")
}

// Creates the binary in the current directory.
func Build() error {
	mg.Deps(Clean)
	mg.Deps(Install)
	mg.Deps(GenerateProto)
	mg.Deps(BuildClient)
	mg.Deps(BuildPublisher)
	return BuildServer()
}

func GenerateProto() error {
	// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto
	return sh.Run("make", "gen-proto")
}

// start client
func StartClient() error {
	mg.Deps(BuildClient)
	return sh.RunV("./bin/grpc-client")
}

// start publisher
func StartPublisher() error {
	mg.Deps(BuildPublisher)
	return sh.RunV("./bin/publisher")
}

// start grpc
func StartGRPC() error {
	mg.Deps(BuildServer)
	return sh.RunV("./bin/grpc-server")
}

// run the test
func Test() error {
	err := sh.RunV("go", "test", "-v", "./...")
	if err != nil {
		return err
	}
	return nil
}
