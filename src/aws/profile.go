package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	ini "gopkg.in/ini.v1"
)

// Profile represents a structure that includes authentication fields necessary
// to authenticate with AWS.
type Profile struct {
	Name            string `ini:"-"`
	AccessKeyID     string `ini:"aws_access_key_id,omitempty"`
	SecretAccessKey string `ini:"aws_secret_access_key,omitempty"`
	SessionToken    string `ini:"aws_session_token,omitempty"`
	MFASerial       string `ini:"mfa_serial,omitempty"`
	Region          string `ini:"region,omitempty"`
}

// ReadConfigFile will read the specific filename and parse the file
// specifically for AWS config file formats and return a list of Profiles.
func ReadConfigFile(filename string) ([]Profile, error) {

	file, err := ini.Load(filename)
	if err != nil {
		return []Profile{}, errors.Wrap(err, "failed to read config file")
	}

	sections := file.SectionStrings()
	var profiles []Profile
	for _, section := range sections {
		// Skip the "DEFAULT" ini section header.
		if section == "DEFAULT" {
			continue
		}
		var profile Profile

		if err = file.Section(section).MapTo(&profile); err != nil {
			return []Profile{}, err
		}

		// Add the section header's name to the Profile.
		profile.Name = section
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// Authenticate will establish a set of temporary AWS credentials for the
// specific duration. The serialNumber and token are linked to the MFA device
// that is connected to the AWS profile/account.
func Authenticate(duration int64, serialNumber, token string) (Profile, error) {

	svc := sts.New(session.New())
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(duration),
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(token),
	}

	result, err := svc.GetSessionToken(input)
	if err != nil {
		return Profile{}, errors.Wrap(err, "get session token failed")
	}

	return Profile{
		AccessKeyID:     *result.Credentials.AccessKeyId,
		SecretAccessKey: *result.Credentials.SecretAccessKey,
		SessionToken:    *result.Credentials.SessionToken,
		MFASerial:       serialNumber,
	}, nil
}
