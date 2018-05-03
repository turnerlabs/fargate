# Fargate CLI

Deploy serverless containers to the cloud from your command line

[![CircleCI](https://circleci.com/gh/turnerlabs/fargate/tree/master.svg?style=svg)](https://circleci.com/gh/turnerlabs/fargate/tree/master)
[![GoDoc](https://godoc.org/github.com/turnerlabs/fargate?status.svg)](https://godoc.org/github.com/turnerlabs/fargate)

## Usage

### Configuration

#### Region

By default, fargate uses *us-east-1* as this is the single region where AWS
Fargate is available. The CLI accepts a --region parameter for future use and
will honor *AWS_REGION* and *AWS_DEFAULT_REGION* environment settings. Note that
specifying a region where all required services aren't available will return an
error.

See the [Region Table][region-table] for a breakdown of what services are
available in which regions.

#### Credentials

fargate is built using the [AWS SDK for Go][go-sdk] which looks for credentials
in the following locations:

1. [Environment Variables][go-env-vars]

1. [Shared Credentials File][go-shared-credentials-file]

1. [EC2 Instance Profile][go-iam-roles-for-ec2-instances]

For more information see [Specifying Credentials][go-specifying-credentials] in
the AWS SDK for Go documentation.

#### Options

There are several ways to specify parameters.  Each item takes precedence over the item below it:

1. CLI arguments (e.g., `--cluster my-cluster`)

1. Environment Variables (e.g., `FARGATE_CLUSTER=my-cluster`)

1. `fargate.yml` (e.g., below)

```yaml
cluster: my-cluster
service: my-service
verbose: false
nocolor: true
```

#### Global Flags

| Flag | Short | Default | Description |
| --- | --- | --- | --- |
| --cluster | -c | | ECS cluster name |
| --region | | us-east-1 | AWS region |
| --no-color | | false | Disable color output |
| --verbose | -v | false | Verbose output |

### Commands

- [Services](#services)

#### Services

Services manage long-lived instances of your containers that are run on AWS
Fargate. If your container exits for any reason, the service scheduler will
restart your containers and ensure your service has the desired number of
tasks running. Services can be used in concert with a load balancer to
distribute traffic amongst the tasks in your service.

- [list](#fargate-service-list)
- [deploy](#fargate-service-deploy)
- [info](#fargate-service-info)
- [logs](#fargate-service-logs)
- [ps](#fargate-service-ps)
- [scale](#fargate-service-scale)
- [env set](#fargate-service-env-set)
- [env unset](#fargate-service-env-unset)
- [env list](#fargate-service-env-list)
- [update](#fargate-service-update)
- [restart](#fargate-service-restart)

##### Flags

| Flag | Short | Default | Description |
| --- | --- | --- | --- |
| --service | -s | | ECS service name |

##### fargate service list

```console
fargate service list
```

List services

##### fargate service deploy

```console
fargate service deploy [--image <docker-image>]
```

Deploy new image to service

The Docker container image to use in the service can be optionally specified
via the --image flag. If not specified, fargate will build a new Docker
container image from the current working directory and push it to Amazon ECR in
a repository named for the task group. If the current working directory is a
git repository, the container image will be tagged with the short ref of the
HEAD commit. If not, a timestamp in the format of YYYYMMDDHHMMSS will be used.

```console
fargate service deploy [--file docker-compose.yml]
```

Deploy image and environment variables defined in a [docker compose file](https://docs.docker.com/compose/overview/) to service

Deploy a docker [image](https://docs.docker.com/compose/compose-file/#image) and [environment variables](https://docs.docker.com/compose/environment-variables/) defined in a docker compose file together as a single unit.  Note that while the docker-compose yaml configuration supports numerous options, only the image and environment variables are deployed to fargate. If the docker compose file defines more than one container, you can use the [label](https://docs.docker.com/compose/compose-file/#labels) `aws.ecs.fargate.deploy: 1` to indicate which container you would like to deploy. For example:

```yaml
version: '3'
services:
  web:
    build: .
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    ports:
    - 80:5000
    environment:
      FOO: bar
      BAZ: bam
    env_file:
    - hidden.env
    labels:
      aws.ecs.fargate.deploy: 1
  redis:
    image: redis
```

##### fargate service info

```console
fargate service info 
```

Inspect service

Show extended information for a service including load balancer configuration,
active deployments, and environment variables.

Deployments show active versions of your service that are running. Multiple
deployments are shown if a service is transitioning due to a deployment or
update to configuration such a CPU, memory, or environment variables.

##### fargate service logs

```console
fargate service logs [--follow] [--start <time-expression>] [--end <time-expression>]
                     [--filter <filter-expression>] [--task <task-id>]
```

Show logs from tasks in a service

Return either a specific segment of service logs or tail logs in real-time
using the --follow option. Logs are prefixed by their log stream name which is
in the format of "fargate/\<service-name>/\<task-id>."

Follow will continue to run and return logs until interrupted by Control-C. If
--follow is passed --end cannot be specified.

Logs can be returned for specific tasks within a service by passing a task ID
via the --task flag. Pass --task with a task ID multiple times in order to
retrieve logs from multiple specific tasks.

A specific window of logs can be requested by passing --start and --end options
with a time expression. The time expression can be either a duration or a
timestamp:

  - Duration (e.g. -1h [one hour ago], -1h10m30s [one hour, ten minutes, and
    thirty seconds ago], 2h [two hours from now])
  - Timestamp with optional timezone in the format of YYYY-MM-DD HH:MM:SS [TZ];
    timezone will default to UTC if omitted (e.g. 2017-12-22 15:10:03 EST)

You can filter logs for specific term by passing a filter expression via the
--filter flag. Pass a single term to search for that term, pass multiple terms
to search for log messages that include all terms. See the [CloudWatch Logs
documentation][cwl-filter-expression] for more details.

##### fargate service ps

```console
fargate service ps
```

List running tasks for a service

##### fargate service scale

```console
fargate service scale <scale-expression>
```

Scale number of tasks in a service

Changes the number of desired tasks to be run in a service by the given scale
expression. A scale expression can either be an absolute number or a delta
specified with a sign such as +5 or -2.

##### fargate service env set

```console
fargate service env set --env <key=value>
```

Set environment variables

At least one environment variable must be specified via the --env flag. Specify
--env with a key=value parameter multiple times to add multiple variables.

##### fargate service env unset

```console
fargate service env unset --key <key-name>
```

Unset environment variables

Unsets the environment variable specified via the --key flag. Specify --key with
a key name multiple times to unset multiple variables.

##### fargate service env list

```console
fargate service env list
```

Show environment variables

##### fargate service update

| Flag | Short | Default | Description |
| --- | --- | --- | --- |
| --cpu | | | Amount of cpu units to allocate for each task |
| --memory | -m | | Amount of MiB to allocate for each task |

```console
fargate service update [--cpu <cpu-units>] [--memory <MiB>]
```

Update service configuration

CPU and memory settings are specified as CPU units and mebibytes respectively
using the --cpu and --memory flags. Every 1024 CPU units is equivilent to a
single vCPU. AWS Fargate only supports certain combinations of CPU and memory
configurations:

| CPU (CPU Units) | Memory (MiB)                          |
| --------------- | ------------------------------------- |
| 256             | 512, 1024, or 2048                    |
| 512             | 1024 through 4096 in 1GiB increments  |
| 1024            | 2048 through 8192 in 1GiB increments  |
| 2048            | 4096 through 16384 in 1GiB increments |
| 4096            | 8192 through 30720 in 1GiB increments |

At least one of --cpu or --memory must be specified.

##### fargate service restart

```console
fargate service restart 
```

Restart service

Creates a new set of tasks for the service and stops the previous tasks. This
is useful if your service needs to reload data cached from an external source,
for example.


[region-table]: https://aws.amazon.com/about-aws/global-infrastructure/regional-product-services/
[go-sdk]: https://aws.amazon.com/documentation/sdk-for-go/
[go-env-vars]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#environment-variables
[go-shared-credentials-file]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#shared-credentials-file
[go-iam-roles-for-ec2-instances]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#iam-roles-for-ec2-instances
[go-specifying-credentials]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials
[cwl-filter-expression]: http://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html#matching-terms-events
