.PHONY: build test clean release-patch release-minor release-major

build:
	go build -o rrk

test:
	go test -v ./...

clean:
	rm -f rrk

patch:
	@./scripts/bump-version.sh patch
	@git push --follow-tags

minor:
	@./scripts/bump-version.sh minor
	@git push --follow-tags

major:
	@./scripts/bump-version.sh major
	@git push --follow-tags

version:
	@git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
