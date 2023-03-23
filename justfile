build:
	goreleaser build --id `go env GOOS` --single-target --snapshot --clean
darwin:
	goreleaser build --id darwin --snapshot --clean
linux:
	goreleaser build --id linux --snapshot --clean
snapshot:
	goreleaser release --snapshot --clean
release:
	git tag `svu next`
	git push --tags
	goreleaser --clean