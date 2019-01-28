// +build mage

// This is the "magefile" for awsctl.
//
// To install mage, run:
// git clone https://github.com/magefile/mage
// cd mage
// go run bootstrap.go
//
// To build awsctl, just run "mage build".

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/magefile/mage/sh"
)

// allow user to override go executable by running as GOEXE=xxx make ... on
// unix-like systems.
var goexe = "go"

var Default = Build

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")
}

// Runs go install for awsctl.
func Build() (err error) {

	ldfTemplate := "--ldflags=%s"

	var ldf string
	env := os.Getenv("ENV")
	if env == "DRYRUN" {
		ldf, err = flags()
		if err != nil {
			return err
		}
	}

	// use -tags make so we can have different behavior for when we know we've
	// built with mage.
	if ldf != "" {
		err = run(goexe, "install", "-mod=vendor", "-tags", "make", fmt.Sprintf(ldfTemplate, ldf), "github.com/outlawlabs/awsctl/cmd/awsctl")
	} else {
		err = run(goexe, "install", "-mod=vendor", "-tags", "make", "github.com/outlawlabs/awsctl/cmd/awsctl")
	}

	return
}

// Runs go mod tidy & vendor.
func Vendor() error {

	log.Println("running go mod tidy")
	if err := sh.Run(goexe, "mod", "tidy"); err != nil {
		return err
	}
	log.Println("running go mod vendor")
	return sh.Run(goexe, "mod", "vendor")
}

// Generates a new release. Expects the TAG environment variable to be set,
// which will create a new tag with that name.
func Release() (err error) {
	releaseTag := regexp.MustCompile(`^v0\.[0-9]+\.[0-9]+$`)
	tag := os.Getenv("TAG")
	if !releaseTag.MatchString(tag) {
		return errors.New("TAG environment variable must be in semver v0.x.x format, but was " + tag)
	}
	if err := sh.RunV("git", "tag", "-a", tag, "-m", tag); err != nil {
		return err
	}
	if err := sh.RunV("git", "push", "origin", tag); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			sh.RunV("git", "tag", "--delete", "$TAG")
			sh.RunV("git", "push", "--delete", "origin", "$TAG")
		}
	}()
	return sh.RunV("goreleaser")
}

// Run go test on code base.
func Test() error {
	return sh.RunV(goexe, "test", "./...")
}

// Run go vet on code base.
func Vet() error {
	return sh.RunV(goexe, "vet", "./...")
}

// Clean up the generated release artifacts.
func Clean() error {
	log.Println("cleaning up dist release directory")
	return sh.RunV("rm", "-rf", "dist")
}
