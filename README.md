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
