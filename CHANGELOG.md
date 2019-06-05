## 0.7.0 (2019-06-05)

### Enhancements

- Propagate tags when creating new task definitions ([#35](https://github.com/turnerlabs/fargate/issues/35))
- Adds describe command that describes a task definition in docker compose format ([#42](https://github.com/turnerlabs/fargate/pull/47))


## 0.6.0 (2019-03-27)

### Enhancements

- Support for secrets in task definitions ([#41](https://github.com/turnerlabs/fargate/pull/41))
- Support for all currently available regions ([#42](https://github.com/turnerlabs/fargate/pull/42))
- Sort env vars output by key ([#33](https://github.com/turnerlabs/fargate/pull/33))


## 0.5.0 (2018-11-14)

### Enhancements

- Only deploy image with docker-compose.yml and --image-only flag ([#23](https://github.com/turnerlabs/fargate/pull/23))
- Support for scheduled tasks - build ([#20](https://github.com/turnerlabs/fargate/issues/20))
- Support for scheduled tasks - deploy ([#25](https://github.com/turnerlabs/fargate/issues/25))
- implements task logs command ([#27](https://github.com/turnerlabs/fargate/pull/27))
- adds --time and --no-prefix to service logs ([#29](https://github.com/turnerlabs/fargate/pull/29))
- Deploy task definition revisions to services ([#28](https://github.com/turnerlabs/fargate/issues/28))

### Bug Fixes

- Fix output of env vars when service.TargetGroupArn is empty ([#22](https://github.com/turnerlabs/fargate/pull/22))


## 0.4.0 (2018-08-22)

### Enhancements

- Add the ability for `service env set` to read environment variables from a file (.e.g.; `fargate service env set -f .env`) ([#17](https://github.com/turnerlabs/fargate/pull/17))
- Add service discovery info to `service info` command ([#21](https://github.com/turnerlabs/fargate/pull/21))


## 0.3.2 (2018-07-05)

### Bug Fixes

- Unset multiple environment variables at once ([#12](https://github.com/turnerlabs/fargate/issues/12))

- Can't add environment variable that contains a comma ([#14](https://github.com/turnerlabs/fargate/issues/14))


## 0.3.1 (2018-06-01)

- Add install link to README

## 0.3.0 (2018-05-16)

### Enhancements

- Console output reworked for consistency and brevity
- macOS users get emoji as a type prefix in console output :tada: -- disable
  with --no-emoji if you're not into fun
- Requests and responses from AWS are displayed in full when --verbose is
  passed
- adds CI/CD pipeline ([#1](https://github.com/turnerlabs/fargate/issues/1))
- updates to turnerlabs and removes unnecessary cmds ([#2](https://github.com/turnerlabs/fargate/issues/2))
- Support more flexible options for configuration (cluster and service) ([#4](https://github.com/turnerlabs/fargate/issues/4))
- Ability to deploy a docker-compose.yml to fargate ([#3](https://github.com/turnerlabs/fargate/issues/3))

### Bug Fixes

- Environment variable service commands now return a polite error message when
  invoked without the service name. ([#22](https://github.com/jpignata/fargate/issues/22)

### Chores

- Utilize `dep` for dependency management
- Add contributor guide, updated license to repo

## 0.2.3 (2018-01-19)

### Features

- Support **--task-role** flag in service create and task run to allow passing
  a role name for the tasks to assume. ([#8][issue-8])

### Enhancements

- Use the `ForceNewDeployment` feature of `UpdateService` in service restart
  instead of incrementing the task definition. ([#14][issue-14])

### Bug Fixes

- Fixed issue where we'd stomp on an existing task role on service updates like
  deployments or environment variable changes. ([#8][issue-8])

## 0.2.2 (2018-01-11)

### Bug Fixes

- Fix service update operation to properly validate and run. ([#11][issue-11])
- Bail out early in service info if the requested service is not active meaning
  it has been previously destroyed.

## 0.2.1 (2018-01-02)

### Bug Fixes

- service create will not run if a load balancer is configured without a port.
- service create and task run will no longer create a repository if an image is
  explictly passed.
- service destroy will remove all references the service's target group and
  delete it.
- Fix git repo detection to properly use a git sha image tag rather than a
  time stamp tag. ([#6][issue-6])
- Fail fast if a user attempts to destroy a service scaled above 0.

## 0.2.0 (2017-12-31)

### Features

- Added **--cluster** global flag to allow running commands against other
  clusters rather than the default. If omitted, the default **fargate** cluster
  is used. ([#2][issue-2])
- lb create, service create, and task run now accept an optional **--subnet-id**
  flag to place resources in different VPCs and subnets rather than the
  defaults. If omitted, resources will be placed within the default subnets
  within the default VPC. ([#2][issue-2])
- lb create, service create, and task run now accept an optional
  **--security-group-id** flag to allow applying more restrictive security
  groups to load balancers, services, and tasks. This flag can be passed
  multiple times to apply multiple security groups. If omitted, a permissive
  security group will be applied.

### Bug Fixes

- Resolved crashes with certificates missing resource records. Certificates that
  fail to be issued immediately after request would cause crashes in lb info and
  lb list as the resource record was never generated.

[issue-2]: https://github.com/jpignata/fargate/issues/2
[issue-6]: https://github.com/jpignata/fargate/issues/6
[issue-8]: https://github.com/jpignata/fargate/issues/8
[issue-11]: https://github.com/jpignata/fargate/issues/11
[issue-14]: https://github.com/jpignata/fargate/issues/14
[issue-22]: https://github.com/jpignata/fargate/issues/22
