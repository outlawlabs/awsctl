package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	ini "gopkg.in/ini.v1"

	"github.com/outlawlabs/awsctl/pkg/logger"
)

const (
	defaultAWSDirectory = "~/.aws"
	credentialsFile     = "~/.aws/credentials"
	configFile          = "~/.aws/config"

	keyAccessKeyID              = "aws_access_key_id"
	keySecretAccessKey          = "aws_secret_access_key"
	keySessionToken             = "aws_session_token"
	keyMFASerial                = "mfa_serial"
	keyRegion                   = "region"
	keyLastAuthentication       = "last_authentication"
	keyAuthenticationExpiration = "authentication_expiration"

	awsCLIHelp = "https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html"

	versionTemplate = `version=%s
timestamp=%s
commit=%s`
)

// Placeholders for binary variable tagging.
var (
	version    = "BREAKING=MINOR++"
	timestamp  = "0 BBY"
	commitHash = string("38762cf7­f55934b3­4d179ae6­a4c80cad­ccbb7f0a")
)

func main() {

	// Set the global ini constants to not require pretty formatting.
	ini.PrettyFormat = false
	ini.PrettyEqual = true

	awsDirectory, err := homedir.Expand(defaultAWSDirectory)
	if err != nil {
		logger.Critical("Failed to expand %s: %s.", defaultAWSDirectory, err)
		os.Exit(1)
	}

	credentialsFile, err := homedir.Expand(credentialsFile)
	if err != nil {
		logger.Critical("Failed to expand ~/.aws/credentials file: %s.", err)
		os.Exit(1)
	}

	configFile, err := homedir.Expand(configFile)
	if err != nil {
		logger.Critical("Failed to expand ~/.aws/config file: %s.", err)
		os.Exit(1)
	}

	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		// Ensure the default ~/.aws directory exists.
		if err = os.MkdirAll(awsDirectory, os.ModePerm); err != nil {
			logger.Critical("Failed to make directory: %s. Error: %s.", awsDirectory, err)
		}

		file, err := os.Create(credentialsFile)
		if err != nil {
			logger.Critical("Failed to create file: %s. Error: %s.", credentialsFile, err)
			os.Exit(1)
		}
		if err = file.Close(); err != nil {
			logger.Warning("Failed to close credentials file.")
		}
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		file, err := os.Create(configFile)
		if err != nil {
			logger.Critical("Failed to create file: %s. Error: %s.", configFile, err)
			os.Exit(1)
		}
		if err = file.Close(); err != nil {
			logger.Warning("Failed to close config file.")
		}
	}

	app := kingpin.New("awsctl", "CLI tool to help manage multiple AWS profiles with MFA enabled.").
		Author("github.com/outlawlabs").
		Version(fmt.Sprintf(versionTemplate, version, timestamp, commitHash))

	configureAuthCommand(app, configFile, credentialsFile)
	configureListCommand(app, configFile, credentialsFile)
	configureNewCommand(app, configFile, credentialsFile)
	configureRemoveCommand(app, configFile, credentialsFile)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

// askForConfirmation asks the user for confirmation. This will not return until
// there is a valid response from the user.
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		logger.Ask(fmt.Sprintf("%s [y/n]: ", s))

		response, err := reader.ReadString('\n')
		if err != nil {
			logger.Critical(err.Error())
			os.Exit(1)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
