APP-BIN := dist/$(shell basename $(shell pwd))
.PHONY: build
build:
	goreleaser build --id $(shell go env GOOS) --single-target --snapshot --clean -o ${APP-BIN}
.PHONY: darwin
darwin:
	goreleaser build --id darwin --snapshot --clean
.PHONY: linux
linux:
	goreleaser build --id linux --snapshot --clean
.PHONY: snapshot
snapshot:
	goreleaser release --snapshot --clean
.PHONY: tag
tag:
	git tag $(shell svu next)
	git push --tags
.PHONY: release
release:
	goreleaser --clean

.PHONY: run
run:
	./${APP-BIN} list
.PHONY: fresh
fresh: build run
.PHONY: lint
lint:
	golangci-lint run -D errcheck
