package main

import (
	"fmt"

	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/outlawlabs/aws-mfa/pkg/aws"
	"github.com/outlawlabs/aws-mfa/pkg/logger"
)

// listCommand represents all of the context for the "list" command.
type listCommand struct {
	configFile string
}

// run will execute the functionality for the "list" command.
func (l *listCommand) run(c *kingpin.ParseContext) error {

	if l.configFile == "" {
		return fmt.Errorf("~/.aws/config file error")
	}

	profiles, err := aws.ReadConfigFile(l.configFile)
	if err != nil {
		return errors.Wrap(err, "failed to read/parse config file")
	}

	if len(profiles) <= 0 {
		logger.Warning("Could not find any profiles. See %s for help.", awsCLIHelp)
		return nil
	}

	headers := []string{"Profile", "MFA Device Serial ARN", "Region"}
	var data [][]string
	for _, value := range profiles {
		data = append(data, []string{value.Name, value.MFASerial, value.Region})
	}
	// Print pretty ASCII table of data.
	logger.Table(headers, data)

	return nil
}

// configureListCommand sets up the "list" command for the main
// kingpin.Application.
func configureListCommand(app *kingpin.Application, configFile, credentialsFile string) {
	c := &listCommand{
		configFile: configFile,
	}
	app.Command("list", "List all AWS profiles.").Action(c.run)
}
