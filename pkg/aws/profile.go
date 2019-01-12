package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	ini "gopkg.in/ini.v1"
)

const (
	keyAccessKeyID     = "aws_access_key_id"
	keySecretAccessKey = "aws_secret_access_key"
	keySessionToken    = "aws_session_token"
	keyMFASerial       = "mfa_serial"
	keyRegion          = "region"
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

// NewProfile will return a new profile based on parameters that are passed in.
func NewProfile(profile, region, serialNumber, accessKeyID, secretAccessKey string) (Profile, error) {

	if profile == "" {
		return Profile{}, errors.New("profile name must be set")
	}
	if region == "" {
		return Profile{}, errors.New("region must be set")
	}
	if err := validateRegion(region); err != nil {
		return Profile{}, err
	}
	if serialNumber == "" {
		return Profile{}, errors.New("serial number must be set")
	}
	if accessKeyID == "" {
		return Profile{}, errors.New("access key must be set")
	}
	if secretAccessKey == "" {
		return Profile{}, errors.New("secret key must be set")
	}
	return Profile{
		Name:            profile,
		Region:          region,
		MFASerial:       serialNumber,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}, nil
}

// validateRegion will examine whether or not the region is supported within
// AWS. The list of regions here is based off of the EC2 supported standard
// public offering.
// See: https://docs.aws.amazon.com/general/latest/gr/rande.html#ec2_region.
func validateRegion(region string) error {
	regions := map[string]bool{
		"us-east-2":      true,
		"us-east-1":      true,
		"us-west-1":      true,
		"us-west-2":      true,
		"ap-south-1":     true,
		"ap-northeast-3": true,
		"ap-northeast-2": true,
		"ap-southeast-1": true,
		"ap-southeast-2": true,
		"ap-northeast-1": true,
		"ca-central-1":   true,
		"cn-north-1":     true,
		"cn-northwest-1": true,
		"eu-central-1":   true,
		"eu-west-1":      true,
		"eu-west-2":      true,
		"eu-west-3":      true,
		"eu-north-1":     true,
		"sa-east-1":      true,
		"us-gov-east-1":  true,
		"us-gov-west-1":  true,
	}
	if ok := regions[region]; !ok {
		return errors.New("region is not supported in AWS. See: https://docs.aws.amazon.com/general/latest/gr/rande.html#ec2_region")
	}
	return nil
}

// Save will persist the Profile's information to their respective places.
func (p Profile) Save(config, credentials *ini.File, configFile, credentialsFile string) error {

	profile := fmt.Sprintf("profile %s", p.Name)

	configSection, err := config.NewSection(profile)
	if err != nil {
		return errors.Wrap(err, "failed to create new config section")
	}
	configSection.Key(keyRegion).SetValue(p.Region)
	configSection.Key(keyMFASerial).SetValue(p.MFASerial)
	if err := config.SaveTo(configFile); err != nil {
		return errors.Wrap(err, "failed to save new config file")
	}

	credentialsSection, err := credentials.NewSection(profile)
	if err != nil {
		return errors.Wrap(err, "failed to create new credentials section")
	}
	credentialsSection.Key(keyAccessKeyID).SetValue(p.AccessKeyID)
	credentialsSection.Key(keySecretAccessKey).SetValue(p.SecretAccessKey)
	if err := credentials.SaveTo(credentialsFile); err != nil {
		return errors.Wrap(err, "failed to save new credentials file")
	}

	return nil
}
