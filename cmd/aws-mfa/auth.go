package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	ini "gopkg.in/ini.v1"

	"github.com/outlawlabs/aws-mfa/src/aws"
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
		return err
	}

	config, err := ini.Load(a.configFile)
	if err != nil {
		return err
	}

	if err := os.Setenv("AWS_SDK_LOAD_CONFIG", "true"); err != nil {
		return err
	}

	if err := os.Setenv("AWS_PROFILE", a.profile); err != nil {
		return err
	}

	configProfile := fmt.Sprintf("profile %s", a.profile)
	if !config.Section(configProfile).HasKey(keyMFASerial) {
		return fmt.Errorf("mfa_serial needs to bet configured for the profile: %s\n", a.profile)
	}
	if !config.Section(configProfile).HasKey(keyRegion) {
		return fmt.Errorf("region needs to bet configured for the profile: %s\n", a.profile)
	}

	mfa := config.Section(configProfile).Key(keyMFASerial).String()
	region := config.Section(configProfile).Key(keyRegion).String()

	prof, err := aws.Authenticate(a.duration, mfa, a.token)
	if err != nil {
		return err
	}

	mfaProfile := fmt.Sprintf("profile %s_mfa", a.profile)
	config.Section(mfaProfile).Key(keyRegion).SetValue(region)
	if err := config.SaveTo(a.configFile); err != nil {
		return err
	}

	mfaProfile = fmt.Sprintf("%s_mfa", a.profile)
	credentials.Section(mfaProfile).Key(keyAccessKeyID).SetValue(prof.AccessKeyID)
	credentials.Section(mfaProfile).Key(keySecretAccessKey).SetValue(prof.SecretAccessKey)
	credentials.Section(mfaProfile).Key(keySessionToken).SetValue(prof.SessionToken)
	credentials.Section(mfaProfile).Key(keyMFASerial).SetValue(prof.MFASerial)
	return credentials.SaveTo(a.credentialsFile)
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
	auth.Flag("profile", "AWS specific profile").Short('p').StringVar(&c.profile)
	auth.Flag("duration", "Active MFA auth duration").Short('d').Int64Var(&c.duration)
}
