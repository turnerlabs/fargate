# Fargate CLI

Deploy serverless containers to the cloud from your command line

[![CircleCI](https://circleci.com/gh/turnerlabs/fargate/tree/master.svg?style=svg)](https://circleci.com/gh/turnerlabs/fargate/tree/master)
[![GoDoc](https://godoc.org/github.com/turnerlabs/fargate?status.svg)](https://godoc.org/github.com/turnerlabs/fargate)

![fargate](fargate.png "fargate")

*fargate* is a command-line interface to deploy containers to [AWS Fargate](https://aws.amazon.com/fargate/). Using *fargate*, developers can easily operate fargate services including things like: deploying applications (images and environment variables), monitoring deployments, viewing container logs, restarting and scaling.

## Install

You can install the latest stable CLI with a curl utility script or by downloading the binary from the releases page. Once installed you'll get the `fargate` command.

```
curl -s get-fargate.turnerlabs.io | sh
```

If you'd like to install the latest prerelease, use this command:

```
curl -s get-fargate.turnerlabs.io | RELEASE=develop sh
```

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
task: my-task
rule: my-event-rule
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
- [Tasks](#tasks)
- [Events](#events)

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

The Docker container image to use in the service can be specified
via the --image flag.


```console
fargate service deploy [--file docker-compose.yml]
```

Deploy image, environment variables, and secrets defined in a [docker compose file](https://docs.docker.com/compose/overview/) to service

Deploy a docker [image](https://docs.docker.com/compose/compose-file/#image) and [environment variables](https://docs.docker.com/compose/environment-variables/) defined in a docker compose file together as a single unit. Note that environments variables and secrets are replaced with what's in the compose file.

Secrets can be defined as key-value pairs under the docker compose file extension field `x-fargate-secrets`. To use extension fields, the compose file version must be  at least `2.4` for the 2.x series or at least `3.7` for the 3.x series.

This allows you to run `docker-compose up` locally to run your app the same way it will run in AWS. Note that while the docker-compose yaml configuration supports numerous options, only the image and environment variables are deployed to fargate. If the docker compose file defines more than one container, you can use the [label](https://docs.docker.com/compose/compose-file/#labels) `aws.ecs.fargate.deploy: 1` to indicate which container you would like to deploy. For example:

```yaml
version: "3.7"
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
    x-fargate-secrets:
      QUX: arn:key:ssm:us-east-1:000000000000:parameter/path/to/my_parameter
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
active deployments, environment variables, and secrets.

Deployments show active versions of your service that are running. Multiple
deployments are shown if a service is transitioning due to a deployment or
update to configuration such a CPU, memory, or environment variables.

##### fargate service logs

```console
fargate service logs [--follow] [--start <time-expression>] [--end <time-expression>]
                     [--filter <filter-expression>] [--task <task-id>]
                     [--time] [--no-prefix]
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

--time includes the log timestamp in the output

--no-prefix excludes the log stream prefix from the output

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
fargate service env set [--env <key=value>] [--file <pathname>]
                        [--secret <key=valueFrom>] [--secret-file <pathname>]
```

Set environment variables and secrets

At least one environment variable or secret must be specified via either the --env,
--file,  --secret, or --secret-file flags. You may specify any number of variables on the command line by
repeating --env before each one, or else place multiple variables in a text
file, one per line, and specify the filename with --file and/or --secret-file.

Each --env and --secret parameter string or line in the file must be of the form
"key=value", with no quotation marks and no whitespace around the "=" unless you want
literal leading whitespace in the value.  Additionally, the "key" side must be
a legal shell identifier, which means it must start with an ASCII letter A-Z or
underscore and consist of only letters, digits, and underscores.

The "value" in "key=value" for each --secret flag should reference the ARN to the AWS Secrets Manager secret or AWS Systems Manager Parameter Store parameter. 

##### fargate service env unset

```console
fargate service env unset --key <key-name>
```

Unset environment variables and secrets

Unsets the environment variable or secret specified via the --key flag. Specify --key with
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


#### Tasks

##### Flags

| Flag | Short | Default | Description |
| --- | --- | --- | --- |
| --task | -t | | Task Definition Family |

Tasks are one-time executions of your container. Instances of your task are run
until you manually stop them either through AWS APIs, the AWS Management
Console, or until they are interrupted for any reason.

- [register](#fargate-task-register)
- [describe](#fargate-task-describe)
- [logs](#fargate-task-logs)


##### fargate task register

```console
fargate task register [--image <docker-image>] 
                      [-e KEY=value -e KEY2=value] [--env-file dev.env]
                      [--secret KEY3=valueFrom] [--secret-file secrets.env]
```

Registers a new [task definition](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html) for the specified docker image, environment variables, or secrets based on the latest revision of the task family and returns the new revision number.

The Docker container image to use in the new Task Definition can be specified
via the --image flag.

The environment variables can be specified using one or many `--env` flags or the `--env-file` flag.

The secrets can be specified using one or many `--secret` flags or the `--secret-file` flag.


```console
fargate task register [--file docker-compose.yml]
```

Registers a new [Task Definition](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html) using the [image](https://docs.docker.com/compose/compose-file/#image), [environment variables](https://docs.docker.com/compose/environment-variables/), and secrets defined in a docker compose file. Note that environments variables are replaced with what's in the compose file.

Secrets can be defined as key-value pairs under the docker compose file extension field `x-fargate-secrets`. To use extension fields, the compose file version must be  at least `2.4` for the 2.x series or at least `3.7` for the 3.x series.

If the docker compose file defines more than one container, you can use the [label](https://docs.docker.com/compose/compose-file/#labels) `aws.ecs.fargate.deploy: 1` to indicate which container you would like to deploy.


##### fargate task describe

```console
fargate task describe
```

The describe command describes a [Task Definition](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html) in [Docker Compose](https://docs.docker.com/compose/overview/) format. The Docker image, environment variables, secrets, and target port are the mapped elements.

This command can be useful for looking at changes made by the `task register`, `service deploy`, or `service env set` commands.  It can also be useful for running a task definition locally for debugging or troubleshooting purposes.

```sh
fargate task describe -t my-app > docker-compose.yml
docker-compose up
```

You can specify the task definition family by using a `fargate.yml` file, the `FARGATE_TASK` envvar, or using 
the `-t` flag, including an optional revision number.

```sh
fargate task describe -t my-app
fargate task describe -t my-app:42
```

Example output:

```yaml
version: "3.7"
services:
  app:
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-app:1.0
    ports:
    - published: 8080
      target: 8080
    environment:
      AWS_REGION: us-east-1
      ENVIRONMENT: dev
      FOO: bar
    x-fargate-secrets:
      KEY: arn:key:ssm:us-east-1:000000000000:parameter/path/to/my_parameter      
    labels:
      aws.ecs.fargate.deploy: "1"
```


##### fargate task logs

```console
fargate task logs [--follow] [--start <time-expression>] [--end <time-expression>]
                  [--filter <filter-expression>] [--task <task-id>] 
                  [--container-name] [--time] [--no-prefix]
```

Show logs from tasks

Assumes a cloudwatch log group with the following convention: `fargate/task/<task>`
where `task` is specified via `--task`, or fargate.yml, or environment variable [options](#options)

Return either a specific segment of task logs or tail logs in real-time using
the --follow option. Logs are prefixed by their log stream name which is in the
format of `fargate/<container-name>/<task-id>.`

`--container-name` allows you to specifiy the container within the task definition to get logs for
(defaults to `app`)

Follow will continue to run and return logs until interrupted by Control-C. If
`--follow` is passed `--end` cannot be specified.

Logs can be returned for specific tasks by passing a task
ID via the `--task` flag. Pass `--task` with a task ID multiple times in order to
retrieve logs from multiple specific tasks.

A specific window of logs can be requested by passing `--start` and `--end` options
with a time expression. The time expression can be either a duration or a
timestamp:

  - Duration (e.g. -1h [one hour ago], -1h10m30s [one hour, ten minutes, and
    thirty seconds ago], 2h [two hours from now])
  - Timestamp with optional timezone in the format of YYYY-MM-DD HH:MM:SS [TZ];
    timezone will default to UTC if omitted (e.g. 2017-12-22 15:10:03 EST)

You can filter logs for specific term by passing a filter expression via the
`--filter` flag. Pass a single term to search for that term, pass multiple terms
to search for log messages that include all terms.

`--time` includes the log timestamp in the output

`--no-prefix` excludes the log stream prefix from the output


#### Events

##### Flags

| Flag | Short | Default | Description |
| --- | --- | --- | --- |
| --rule | -r | | CloudWatch Events Rule |

The `events` command provides subcommands for working with [CloudWatch Events](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/WhatIsCloudWatchEvents.html) (scheduled tasks, etc.)

- [target](#fargate-events-target)


##### fargate events target

```console
fargate events target --revision <revision>
```

"Deploys" (causes the next event rule invocation to run the new version) a task definition revision to a CloudWatch Event Rule by updating the rule target's `EcsParameters.TaskDefinitionArn`.

A typical CI/CD system might do something like:
```console
REVISION=$(fargate task register -i 123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:${VERSION}-${CIRCLE_BUILD_NUM} -e FOO=bar)
fargate events target -r ${REVISION}
```


[region-table]: https://aws.amazon.com/about-aws/global-infrastructure/regional-product-services/
[go-sdk]: https://aws.amazon.com/documentation/sdk-for-go/
[go-env-vars]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#environment-variables
[go-shared-credentials-file]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#shared-credentials-file
[go-iam-roles-for-ec2-instances]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#iam-roles-for-ec2-instances
[go-specifying-credentials]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials
[cwl-filter-expression]: http://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html#matching-terms-events
