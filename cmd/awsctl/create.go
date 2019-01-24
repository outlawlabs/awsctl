package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	ini "gopkg.in/ini.v1"

	"github.com/outlawlabs/awsctl/pkg/aws"
	"github.com/outlawlabs/awsctl/pkg/logger"
)

// newCommand represents all of the context for the "new" command.
type newCommand struct {
	profile         string
	region          string
	accessKey       string
	secretKey       string
	mfaSerial       string
	configFile      string
	credentialsFile string
}

// run will execute the functionality for the "new" command.
func (n *newCommand) run(c *kingpin.ParseContext) error {

	if n.configFile == "" {
		return errors.New("~/.aws/config file error")
	}

	if n.credentialsFile == "" {
		return errors.New("~/.aws/credentials file error")
	}

	credentials, err := ini.Load(n.credentialsFile)
	if err != nil {
		return errors.Wrap(err, "failed to read credentials file")
	}

	config, err := ini.Load(n.configFile)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	// Create the template for the new AWS CLI profile.
	configProfile := fmt.Sprintf("profile %s", n.profile)

	// Ignore error -- *should* only come back as not found if the section does
	// not exist already.
	sections := config.SectionStrings()
	for i := 0; i < len(sections); i++ {
		if sections[i] == configProfile {
			return errors.New("cannot create new profile, it already exists")
		}
	}
	// Also check for the same section header within the credentials file.
	sections = credentials.SectionStrings()
	for i := 0; i < len(sections); i++ {
		if sections[i] == configProfile {
			return errors.New("cannot create new profile, it already exists")
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	logger.Ask("Enter the AWS region you want to save:")
	scanner.Scan()
	region := scanner.Text()

	logger.Ask("Enter your MFA serial number for your IAM user:")
	scanner.Scan()
	mfaSerial := scanner.Text()

	logger.Ask("Enter your generated access key ID:")
	scanner.Scan()
	accessKeyID := scanner.Text()

	logger.Ask("Enter your generated secret access key:")
	scanner.Scan()
	secretAccessKey := scanner.Text()

	fmt.Println("")
	logger.Info("Working on your new AWS profile: %s", n.profile)
	profile, err := aws.NewProfile(n.profile, region, mfaSerial, accessKeyID, secretAccessKey)
	if err != nil {
		return err
	}

	// Save the new profile to respective files.
	if err := profile.Save(config, credentials, n.configFile, n.credentialsFile); err != nil {
		return err
	}

	logger.Success("Successfully saved new config and credentials for profile: %s.", n.profile)
	logger.Always("Start using your new profile: aws-mfa auth --help")
	return nil
}

// configureNewCommand sets up the "new" command for the main
// kingpin.Application.
func configureNewCommand(app *kingpin.Application, configFile, credentialsFile string) {
	n := &newCommand{
		configFile:      configFile,
		credentialsFile: credentialsFile,
	}
	new := app.Command("new", "Save a new AWS profile & credential pair.").Action(n.run)
	new.Arg("profile", "AWS profile to create.").Required().StringVar(&n.profile)
}
