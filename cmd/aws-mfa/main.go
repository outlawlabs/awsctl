package main

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	ini "gopkg.in/ini.v1"

	"github.com/outlawlabs/aws-mfa/src/logger"
)

const (
	credentialsFile = "~/.aws/credentials"
	configFile      = "~/.aws/config"

	keyAccessKeyID     = "aws_access_key_id"
	keySecretAccessKey = "aws_secret_access_key"
	keySessionToken    = "aws_session_token"
	keyMFASerial       = "mfa_serial"
	keyRegion          = "region"
)

var (
	version    = "DEV"
	timestamp  = ""
	commitHash = ""
)

func main() {

	// Set the global ini constants to not require pretty formatting.
	ini.PrettyFormat = false
	ini.PrettyEqual = true

	credentialsFile, err := homedir.Expand("~/.aws/credentials")
	if err != nil {
		logger.Critical("Failed to expand ~/.aws/credentials file: %s.", err)
		os.Exit(1)
	}

	configFile, err := homedir.Expand("~/.aws/config")
	if err != nil {
		logger.Critical("Failed to expand ~/.aws/config file: %s.", err)
		os.Exit(1)
	}

	app := kingpin.New("aws-mfa", "CLI tool to help manage multiple AWS profiles with MFA requirements.").
		Author("github.com/outlawlabs").
		Version(version)
	configureAuthCommand(app, configFile, credentialsFile)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
