# looking-glass

Looking Glass is a tool for mirroring objects to Artifactory.  Currently only S3 is supported as a source, but FTP support will be added in the near future.

## Why would you want to do this?

It is common for vendors to release binaries to an S3 bucket or FTP.  To prevent 100's of servers from reaching outside our network during deployments, we had manually copied these binaries to our Artifactory instance.  Looking Glass automates this for us, so we know our Artifactory repo is always up to date! 


# Setup

The latest version of looking-glass can be found on the [Releases](https://github.com/simplifi/looking-glass/releases) tab.

First, you'll need to create a "Generic" repository in Artifactory, and ensure you have a user that can write to it.

Second, you'll need to gather up the credentials for the S3 bucket you'd like to mirror (you'll need them for the config in the next step)

Third, create a looking-glass yaml configuration file that tells it how to talk to Artifactory/S3, and the source/destination of the objects (see details below)

### Example Configuration:
```yaml
artifactory:
  url: http://my.artifactory.server/artifactory/
  username: my-artifactory-user
  key: my-artifactory-key
agents:
  - name: my-agent-name
    aws_bucket: my-s3-bucket
    aws_key: my-aws-key
    aws_secret: my-aws-secret
    aws_prefix: my-prefix
    aws_region: us-west-2
    artifactory_repo: my-repo
    sleep_duration: 900
```

### `artifactory`
This is where you tell looking-glass how to talk to your Artifactory server
- `url` - The URL to your Artifactory server
- `username` - The username to use when authenticating with Artifactory
- `key` - The user's key used when authenticating with Artifactory

### `agents`
This is where you tell looking-glass about the agent(s) configuration
- `name` - The name of this agent, mainly used in logging
- `aws_bucket` - The bucket from which you wish to mirror
- `aws_key` - The AWS Key ID to use when authenticating with S3
- `aws_secret` - The AWS Secret Key to use when authenticating with S3
- `aws_prefix` - The prefix to mirror from the S3 bucket
- `aws_region` - The region in which the S3 bucket exists
- `artifactory_repo` - The name of the Artifactory repo which will be the destination for the mirrored objects
- `sleep_duration` - How long to wait before polling the S3 bucket for changes (in seconds)

# Usage

### Basic Usage
```shell script
Looking Glass (S3->Artifactory Mirror)

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


