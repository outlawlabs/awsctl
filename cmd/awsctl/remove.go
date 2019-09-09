package main

import (
	"fmt"

	"github.com/outlawlabs/awsctl/pkg/logger"
	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	ini "gopkg.in/ini.v1"

	"github.com/outlawlabs/awsctl/pkg/aws"
)

// removeCommand represents all of the context for the "remove" command.
type removeCommand struct {
	profile         string
	configFile      string
	credentialsFile string
}

// run will execute the functionality for the "new" command.
func (r *removeCommand) run(c *kingpin.ParseContext) error {

	if r.configFile == "" {
		return errors.New("~/.aws/config file error")
	}

	if r.credentialsFile == "" {
		return errors.New("~/.aws/credentials file error")
	}

	credentials, err := ini.Load(r.credentialsFile)
	if err != nil {
		return errors.Wrap(err, "failed to read credentials file")
	}

	config, err := ini.Load(r.configFile)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	// Create the template for the new AWS CLI profile.
	configProfile := fmt.Sprintf("profile %s", r.profile)

	// Ignore error -- *should* only come back as not found if the section does
	// not exist already.
	if _, err = config.GetSection(configProfile); err != nil {
		return errors.New("cannot remove profile, it does not exist")
	}

	confirm := askForConfirmation(fmt.Sprintf("Are you sure you want to remove your profile: %s", r.profile))
	if !confirm {
		logger.Info("Cancelled removal.")
		return nil
	}

	if err = aws.RemoveProfile(configProfile, config, credentials, r.configFile, r.credentialsFile); err != nil {
		return err
	}

	logger.Success("Successfully removed config and credentials for profile: %s.", r.profile)
	return nil
}

// configureRemoveCommand sets up the "remove" command for the main
// kingpin.Application.
func configureRemoveCommand(app *kingpin.Application, configFile, credentialsFile string) {
	r := &removeCommand{
		configFile:      configFile,
		credentialsFile: credentialsFile,
	}
	remove := app.Command("remove", "Remove an existing AWS profile & credential pair.").Action(r.run)
	remove.Arg("profile", "AWS profile to remove.").Required().StringVar(&r.profile)
}
