# looking-glass

## Looking glass is no longer actively maintained.

[![Build Status](https://travis-ci.com/simplifi/looking-glass.svg?branch=master)](https://travis-ci.com/simplifi/looking-glass) [![Go Report Card](https://goreportcard.com/badge/github.com/simplifi/looking-glass)](https://goreportcard.com/report/github.com/simplifi/looking-glass) [![Release](https://img.shields.io/github/release/simplifi/looking-glass.svg)](https://github.com/simplifi/looking-glass/releases/latest)

Looking Glass is a tool for mirroring objects to Artifactory.

## Why would you want to do this?

It is common for vendors to release binaries to an S3 bucket or FTP.  To prevent 100's of servers from reaching outside our network during deployments, we had manually copied these binaries to our Artifactory instance.  Looking Glass automates this for us, so we know our Artifactory repo is always up to date!


# Setup

The latest version of looking-glass can be found on the [Releases](https://github.com/simplifi/looking-glass/releases) tab.

First, you'll need to create a "Generic" repository in Artifactory, and ensure you have a user that can write to it.

Second, you'll need to gather up the credentials for the source you'd like to mirror (you'll need them for the config in the next step)

Third, create a looking-glass yaml configuration file that tells it how to talk to the source/destination of the objects (see details below)

## Example Configuration:
```yaml
artifactory:
  url: http://my.artifactory.server/artifactory/
  username: my-artifactory-user
  key: my-artifactory-key
agents:
  - name: my-s3-agent
    artifactory_repo: my-repo-s3
    sleep_duration: 900
    downloader:
      type: s3
      config:
        aws_bucket: my-s3-bucket
        aws_key: my-aws-key
        aws_secret: my-aws-secret
        aws_prefix: my-prefix
        aws_region: us-west-2
  - name: my-github-agent
    artifactory_repo: my-repo-github
    sleep_duration: 900
    downloader:
      type: github
      config:
        github_repo: simplifi/looking-glass
        github_token: my-github-token
```

### `artifactory`
This is where you tell looking-glass how to talk to your Artifactory server
- `url` - The URL to your Artifactory server
- `username` - The username to use when authenticating with Artifactory
- `key` - The user's key used when authenticating with Artifactory

### `agents`
This is where you tell looking-glass about the agent(s) configuration
- `name` - The name of this agent, mainly used in logging
- `artifactory_repo` - The name of the Artifactory repo which will be the destination for the mirrored objects
- `sleep_duration` - How long to wait before polling the for changes (in seconds)

### `agents.downloader` (s3)
This is where you tell looking-glass how to download objects from s3
- `type` -  The type of downloader that you with to run (`s3` in this case)
- `config.aws_bucket` - The bucket from which you wish to mirror
- `config.aws_key` - The AWS Key ID to use when authenticating with S3
- `config.aws_secret` - The AWS Secret Key to use when authenticating with S3
- `config.aws_prefix` - The prefix to mirror from the S3 bucket
- `config.aws_region` - The region in which the S3 bucket exists

### `agents.downloader` (github)
This is where you tell looking-glass how to download assets from Github
- `type` -  The type of downloader that you with to run (`github` in this case)
- `config.github_repo` - The github repo (in the form of `owner/repo_name`) from which to pull release assets
- `config.github_token` - (optional) The token to authenticate with when pulling release assets

# Usage

### Basic Usage
```
Looking Glass (Artifactory Mirror)

Usage:
  looking-glass [command]

Available Commands:
  help        Help about any command
  start       Start the Looking Glass agent
  version     Print the version number of looking-glass

Flags:
  -h, --help   help for looking-glass

Use "looking-glass [command] --help" for more information about a command.
```

### To start up the agents:
```shell script
looking-glass start -c /path/to/your/config.yml
```

# Development

### Compiling
```shell script
make build
```

### Running Tests
To run all the standard tests:
```shell script
make test
```

### Releasing
This project is using [goreleaser](https://goreleaser.com). GitHub release creation is automated using Travis CI. New releases are automatically created when new tags are pushed to the repo.
```shell script
$ TAG=0.1.0 make tag
```
