package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	ini "gopkg.in/ini.v1"

	"github.com/outlawlabs/awsctl/pkg/aws"
	"github.com/outlawlabs/awsctl/pkg/logger"
)

// authCommand represents all of the context for the "auth" command.
type authCommand struct {
	token           string
	profile         string
	duration        int64
	configFile      string
	credentialsFile string
}

// run will execute the functionality for the "auth" command.
func (a *authCommand) run(c *kingpin.ParseContext) error {

	if a.configFile == "" {
		return fmt.Errorf("~/.aws/config file error")
	}

	if a.credentialsFile == "" {
		return fmt.Errorf("~/.aws/credentials file error")
	}

	credentials, err := ini.Load(a.credentialsFile)
	if err != nil {
		return errors.Wrap(err, "failed to read credentials file")
	}

	config, err := ini.Load(a.configFile)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	if err := os.Setenv("AWS_SDK_LOAD_CONFIG", "true"); err != nil {
		return errors.Wrap(err, "failed to set AWS_SDK_LOAD_CONFIG value")
	}

	if err := os.Setenv("AWS_PROFILE", a.profile); err != nil {
		return errors.Wrap(err, "failed to set AWS_PROFILE value")
	}

	configProfile := fmt.Sprintf("profile %s", a.profile)
	if !config.Section(configProfile).HasKey(keyMFASerial) {
		return fmt.Errorf("mfa_serial needs to bet configured for the profile: %s", a.profile)
	}
	if !config.Section(configProfile).HasKey(keyRegion) {
		return fmt.Errorf("region needs to bet configured for the profile: %s", a.profile)
	}

	mfa := config.Section(configProfile).Key(keyMFASerial).String()
	region := config.Section(configProfile).Key(keyRegion).String()

	logger.Info("Attempting to authenticate with credentials for profile: %s.", a.profile)
	prof, err := aws.Authenticate(a.duration, mfa, a.token)
	if err != nil {
		return err
	}

	mfaProfile := fmt.Sprintf("profile %s_mfa", a.profile)
	config.Section(mfaProfile).Key(keyRegion).SetValue(region)
	if err := config.SaveTo(a.configFile); err != nil {
		return errors.Wrap(err, "failed to save new config file")
	}

	mfaProfile = fmt.Sprintf("%s_mfa", a.profile)
	credentials.Section(mfaProfile).Key(keyAccessKeyID).SetValue(prof.AccessKeyID)
	credentials.Section(mfaProfile).Key(keySecretAccessKey).SetValue(prof.SecretAccessKey)
	credentials.Section(mfaProfile).Key(keySessionToken).SetValue(prof.SessionToken)
	credentials.Section(mfaProfile).Key(keyMFASerial).SetValue(prof.MFASerial)
	if err := credentials.SaveTo(a.credentialsFile); err != nil {
		return errors.Wrap(err, "failed to save new credentials file")
	}

	logger.Success("Successfully created a MFA authenticated session for profile: %s.", a.profile)
	logger.Always("Activate your MFA profile: export AWS_PROFILE=%s_mfa", a.profile)
	return nil
}

// configureAuthCommand sets up the "auth" command for the main
// kingpin.Application.
func configureAuthCommand(app *kingpin.Application, configFile, credentialsFile string) {
	c := &authCommand{
		configFile:      configFile,
		credentialsFile: credentialsFile,
	}
	auth := app.Command("auth", "MFA authentication.").Action(c.run)
	auth.Flag("token", "One time MFA token.").Short('t').StringVar(&c.token)
	auth.Flag("profile", "AWS specific profile.").Short('p').StringVar(&c.profile)
	auth.Flag("duration", "Active MFA auth duration.").Short('d').Int64Var(&c.duration)
}
