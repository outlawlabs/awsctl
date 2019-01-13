// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func flags() (string, error) {
	timestamp := time.Now().Format(time.RFC3339)
	hash, err := output("git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	version, err := gitTag()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`-X "github.com/outlawlabs/cmd/awsctl.timestamp=%s" -X "github.com/outlawlabs/cmd/awsctl.commitHash=%s" -X "github.com/outlawlabs/cmd/awsctl.version=%s"`, timestamp, hash, version), nil
}

func gitTag() (string, error) {
	s, err := output("git", "describe", "--tags")
	if err != nil {
		ee, ok := errors.Cause(err).(*exec.ExitError)
		if ok && ee.Exited() {
			return "dev", nil
		}
		return "", err
	}

	return strings.TrimSuffix(s, "\n"), nil
}

func rm(s string) error {
	err := os.RemoveAll(s)
	if os.IsNotExist(err) {
		return nil
	}
	return errors.Wrapf(err, `failed to remove %s`, s)
}

func run(cmd string, args ...string) error {
	return runWith(nil, cmd, args...)
}

func runWith(env []string, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	for _, v := range env {
		c.Env = append(c.Env, v)
	}
	c.Stderr = os.Stderr
	if os.Getenv("MAGEFILE_VERBOSE") != "" {
		c.Stdout = os.Stdout
	}
	return errors.Wrapf(c.Run(), `failed to run %v %q`, cmd, args)
}

func output(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	c.Stderr = os.Stderr
	b, err := c.Output()
	if err != nil {
		return "", errors.Wrapf(err, `failed to run %v %q`, cmd, args)
	}
	return string(b), nil
}
