package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Profile represents a structure that includes authentication fields necessary
// to authenticate with AWS.
type Profile struct {
	AccessKeyID     string `ini:"aws_access_key_id,omitempty"`
	SecretAccessKey string `ini:"aws_secret_access_key,omitempty"`
	SessionToken    string `ini:"aws_session_token,omitempty"`
	MFASerial       string `ini:"mfa_serial,omitempty"`
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
		return Profile{}, err
	}

	return Profile{
		AccessKeyID:     *result.Credentials.AccessKeyId,
		SecretAccessKey: *result.Credentials.SecretAccessKey,
		SessionToken:    *result.Credentials.SessionToken,
		MFASerial:       serialNumber,
	}, nil
}
