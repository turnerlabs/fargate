# How to contribute

Thanks for your interest in the project! We want to welcome contributors so we put together the following set of guidelines to help participate.

## Workflow

- **Did you find a bug?**

  Awesome! Please feel free to open an issue first, or if you have a fix open a
  pull request that describes the bug with code that demonstrates the bug in a
  test and addresses it.

- **Do you want to add a feature?**

  Features begin life as a proposal. Please open a pull request with a proposal
  that explains the feature, its use case, considerations, and design. This will
  allow interested contributors to weigh in, refine the idea, and ensure there's
  no wasted time in the event a feature doesn't fit with our direction.

## Setup

- Ensure you're using the Go version specified in [go.mod](go.mod). Check the required version:

  ```console
  go version
  ```

- Clone the repository

  ```console
  git clone https://github.com/turnerlabs/fargate.git
  cd fargate
  ```

- Download dependencies using Go modules

  ```console
  go mod download
  ```

- Make sure you can run the tests

  ```console
  make test
  ```

## Testing

- Tests can be run via `go test` or `make test`

- To generate mocks as you add functionality, run `make mocks` or use `go
generate` directly

## Building

- To build a binary for your platform run `make`

- For cross-building for all supported platforms, run `make dist` which builds
  binaries for all supported platforms.

## Making Changes

- Create a feature branch from where you want to base your work.
  - This is usually the `master` branch.
  - To quickly create a feature branch; `git checkout -b feature/my-feature`. Please avoid working directly on the
    `master` branch.
- Make commits of logical units.
- Run `go fmt ./cmd` before committing.
- Make sure you have added the necessary tests for your changes.
- Run _all_ the tests to assure nothing else was accidentally broken.

## Submitting Changes

- Push your changes to a feature branch in your fork of the repository.
- Submit a pull request to the `master` branch to the repository in the turnerlabs organization.

## Release Process

- After a feature pull request has been merged into the `master` branch, a maintainer will tag it with a pre-release number which will trigger the creation of a release.
- The release branch is merged to `master`, tagged, and pushed (along with tags).
- After sufficient testing, this feature will be included in a full release with a normal version tag.

## Licensing

This project is released under the [Apache 2.0 license][apache].
