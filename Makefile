APP-BIN := dist/$(shell basename $(shell pwd))

.PHONY: build darwin fresh lint linux qa release run snapshot tag test watch
build:
	goreleaser build --id $(shell go env GOOS) --single-target --snapshot --clean -o ${APP-BIN}
darwin:
	goreleaser build --id darwin --snapshot --clean
linux:
	goreleaser build --id linux --snapshot --clean
snapshot:
	goreleaser release --snapshot --clean
tag:
	git tag $(shell svu next)
	git push --tags
release: tag
	goreleaser --clean

watch:
	gotestsum --watch --format testname
lint:
	pre-commit run --files $(shell git ls-files -m)
test:
	gotestsum --format testname
qa: lint test

run:
	./${APP-BIN} list
fresh: build run
