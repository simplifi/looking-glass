TAG?=""

# Clean up any cruft left over from old builds
.PHONY: clean
clean:
	rm -rf looking-glass dist/

# Build a beta version of looking-glass
.PHONY: build
build: clean
	CGO_ENABLED=0 go build

# Run all tests
.PHONY: test
test: fmt lint vet test-unit

# Run a test release with goreleaser
.PHONY: test-release
test-release:
	goreleaser --snapshot --skip-publish --rm-dist

# Run unit tests
.PHONY: test-unit
test-unit:
	go test -v -race ./...

# Check formatting
.PHONY: fmt
fmt:
	test -z "$(shell gofmt -l .)"

# Run linter
.PHONY: lint
lint:
	golint -set_exit_status ./...

# Run vet
.PHONY: vet
vet:
	go vet ./...

# For use in ci
.PHONY: ci
ci: build test

# Create a git tag
.PHONY: tag
tag:
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)

# Requires GITHUB_TOKEN environment variable to be set
.PHONY: release
release: clean
	goreleaser
