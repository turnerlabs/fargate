.PHONY: mocks test build dist

PACKAGES := $(shell go list ./... | grep -v /mock)

mocks:
	go get github.com/golang/mock/mockgen
	go generate $(PACKAGES)

test:
	go test -race -cover $(PACKAGES)

build:
	go build -o bin/fargate main.go

dist:
	GOOS=darwin GOARCH=amd64 go build -o dist/build/fargate-darwin-amd64/fargate main.go
	GOOS=linux GOARCH=amd64 go build -o dist/build/fargate-linux-amd64/fargate main.go
	GOOS=linux GOARCH=386 go build -o dist/build/fargate-linux-386/fargate main.go
	GOOS=linux GOARCH=arm go build -o dist/build/fargate-linux-arm/fargate main.go

	cd dist/build/fargate-darwin-amd64 && zip fargate-${FARGATE_VERSION}-darwin-amd64.zip fargate
	cd dist/build/fargate-linux-amd64 && zip fargate-${FARGATE_VERSION}-linux-amd64.zip fargate
	cd dist/build/fargate-linux-386  && zip fargate-${FARGATE_VERSION}-linux-386.zip fargate
	cd dist/build/fargate-linux-arm  && zip fargate-${FARGATE_VERSION}-linux-arm.zip fargate

	find dist/build -name *.zip -exec mv {} dist \;

	rm -rf dist/build

xplat-build:
	BUILD_VERSION=$(git describe --tags)
	echo building ${BUILD_VERSION}
	gox -osarch="darwin/amd64" -osarch="linux/386" -osarch="linux/amd64" -osarch="windows/amd64" \
		-ldflags "-X main.version=${BUILD_VERSION}" -output "dist/ncd_{{.OS}}_{{.Arch}}"

prerelease:
	ghr --prerelease -u ${CIRCLE_PROJECT_USERNAME} -r ${CIRCLE_PROJECT_REPONAME} --replace `git describe --tags` dist/

release:
	ghr -u ${CIRCLE_PROJECT_USERNAME} -r ${CIRCLE_PROJECT_REPONAME} --replace `git describe --tags` dist/