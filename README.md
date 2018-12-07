# ghuser

[![Go Report Card](https://goreportcard.com/badge/github.com/niavasha/ghuser)](https://goreportcard.com/report/github.com/niavasha/ghuser) <a href="https://godoc.user/github.com/niavasha/ghuser"><img src="https://godoc.user/github.com/niavasha/ghuser?status.svg" alt="GoDoc"></a> [![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.user/licenses/Apache-2.0)

ghuser allows you to quickly clone all of an users repos into a single directory. This can be useful in many situations including

1. Searching your users codebase with ack, silver searcher, grep etc..
2. Bash scripting
3. Creating backups
4. Onboarding new teammates
5. Performing Audits

ALL CREDITS ARE OWED TO gabrie30 who wrote ghorg - I simply refactored it so it could be used for users rather than organizations. You can find gabrie30's originaly code here: https://github.com/gabrie30/ghorg - Props due!

> When running ghuser a second time, all local changes in your *_ghuser directory will be overwritten by whats on GitHub. If you are working out of this directory, make sure you rename it before running a second time otherwise all of you changes will be lost.

## Setup

### Homebrew

> optional

```bash
$ brew update
$ brew upgrade git
```
> required

```bash
$ brew install niavasha/utils/ghuser
$ curl https://raw.githubusercontent.com/niavasha/ghuser/master/.env-sample > $HOME/.ghuser
```

### Golang

```bash
$ go get -u github.com/niavasha/ghuser
$ cd $HOME/go/src/github.com/niavasha/ghuser
$ cp .env-sample .env
# update your .env, if needed
# If ghuser_GITHUB_TOKEN is not set in .ghuser, defaults to keychain, see below
$ make install
$ go install
```

## Use

```bash
$ ghuser user-you-want-to-clone
```

> ghuser defaults to master however, for gitflows you can run on develop by setting ghuser_BRANCH=develop or similar

## Configuration

All configuration will be done in the .ghuser file. This file will be created from the [.env-sample](https://github.com/niavasha/ghuser/blob/master/.env-sample) and copied into `~/.ghuser`. Make sure this file exists then configure to your needs.

## Default GitHub Token Used

```bash
$ security find-internet-password -s github.com  | grep "acct" | awk -F\" '{ print $4 }'
```

> If running this does not return the correct key you will need to generate a token via GithHub and add it to your $HOME/.ghuser

> To view all other default environment variables see .env-sample

## Auth through SSO

- If user is behind SSO a normal token will not work. You will need to add SSO to the [Github token](https://help.github.com/articles/authorizing-a-personal-access-token-for-use-with-a-saml-single-sign-on-useranization/)

## Troubleshooting

- Make sure your `$ git --version` is >= 2.19.0
- You may need to increase your ulimits if cloning a large user
- Other issues can most likely be resolved by adding a `.ghuser` to your users home directory and setting the necessary values defined in the `.env-sample`

### Updating brew tap
- [See Readme](https://github.com/niavasha/homebrew-utils/blob/master/README.md)
