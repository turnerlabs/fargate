# How to contribute

Thanks for your interest in the project!  We want to welcome contributors so we put together the following set of guidelines to help participate.

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

- Ensure you're using golang 1.10+.

  ```console
  go version
  ```

- Install [`dep`][dep] if not present on your system. See their [installation
  instructions][dep-install] and [releases page][dep-releases] for details.

- Install the source code from GitHub

  ```console
  go get github.com/turnerlabs/fargate
  ```

- Run `dep ensure` to install required dependencies

  ```console
  cd $GOPATH/src/github.com/turnerlabs/fargate
  dep ensure
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
  binaries for darwin (64-bit) and linux (Arm, 32-bit, 64-bit).

## Making Changes

* Create a feature branch from where you want to base your work.
  * This is usually the `develop` branch.
  * To quickly create a feature branch; `git checkout -b feature/my-feature`. Please avoid working directly on the
    `master` branch.
* Make commits of logical units.
* Run `go fmt ./cmd` before committing.
* Make sure you have added the necessary tests for your changes.
* Run _all_ the tests to assure nothing else was accidentally broken.

## Submitting Changes

* Push your changes to a feature branch in your fork of the repository.
* Submit a pull request to the `develop` branch to the repository in the turnerlabs organization.

## Release Process

* After a feature pull request has been merged into the `develop` branch, a CI build will be automatically kicked off.  The CI build will run unit tests, do a multi-platform build and automatically deploy the build to the [Github releases](https://github.com/turnerlabs/fargate/releases) page as a pre-release using the latest tag (`git describe --tags`) as the version number.
* After the core team decides which features will be included in the next release, a release branch is created (e.g., `release/v0.5`) from develop.
* The `CHANGELOG.md` file is updated to document the release in the release branch.
* The release branch is merged to `master`, tagged, and pushed (along with tags).
* This will kick off a build that builds using the latest tag and deploys as a Github release.
* The release branch is then merged back to `develop`, tagged for pre-release (to start next version, e.g. v0.6.0-pre) and pushed.

## Licensing

This project is released under the [Apache 2.0 license][apache].
