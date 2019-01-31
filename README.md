# awsctl

[![Outlaw Labs](https://img.shields.io/badge/developed%20by-Outlaw%20Labs-%23ff9933.svg)](https://github.com/outlawlabs)
[![Built with Go](https://img.shields.io/badge/built%20with-Go-blue.svg)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/outlawlabs/awsctl)](https://goreportcard.com/badge/github.com/outlawlabs/awsctl)
[![MIT License](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php)
[![CircleCI](https://circleci.com/gh/outlawlabs/awsctl.svg?style=svg)](https://circleci.com/gh/outlawlabs/awsctl)

### CLI based tool to help manage AWS profiles for account enabled with MFA.

## Purpose

One main problem that teams face when enforcing MFA device authentication while
working with AWS CLI profiles is there are no official stream-lined tool, or
tools, to manage your temporary credential sessions easily.

Instead of piecing together some _bash_, or _shell_, script to manage and do the
magic behind the scenes to authenticate with your credentials, then get the new
temporary credentials and save them to your existing file, create a new
config/credentials file, or however you might your flow look. This tool is
designed to be a **single binary** that will enable you to **create**, **list**,
**authenticate** and generally **manage** your AWS profiles.

As a side note -- to understand the general flow of how to do this natively with
the AWS CLI check out this AWS support [article](https://aws.amazon.com/premiumsupport/knowledge-center/authenticate-mfa-cli/).

## Install

To install `awsctl` checkout the [releases](https://github.com/outlawlabs/awsctl/releases)
page to find the latest download for your operating system.

## Usage

To get started with `awsctl` you will need to ensure you have an AWS IAM account
already setup with an MFA device. For more information about checking your MFA
status checkout the AWS [docs](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_checking-status.html).

### New Profile

To create a new MFA profile with `awsctl` you simply grab your access key
credentials from your account, as well as your _Assigned MFA device_.

You will need to use these when you generate a new AWS profile with `awsctl`.

Example interactive creation process --

```sh
$ awsctl new example
[?]  Enter the AWS region you want to save:
us-east-1
[?]  Enter your MFA serial number for your IAM user:
arn:aws:iam::123456789012:mfa/cowboy
[?]  Enter your generated access key ID:
********************
[?]  Enter your generated secret access key:
****************************************

[ℹ]  Working on your new AWS profile: example
[✔]  Successfully saved new config and credentials for profile: example.
[✈]  Start using your new profile: awsctl auth --help
```

### Authenticate

When you need to authenticate and create a new temporary session for our AWS CLI
interactions. We leverage the streamlined functionality in `awsctl new` command.

Example authentication process --

```sh
$ awsctl auth --profile example --duration 129000 --token 639959
[ℹ]  Attempting to authenticate with credentials for profile: example.
[✔]  Successfully created a MFA authenticated session for profile: example.
[✈]  Activate your MFA profile: export AWS_PROFILE=example_mfa
```
