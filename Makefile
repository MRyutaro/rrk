.PHONY: build test clean release-patch release-minor release-major

build:
	go build -o rrk

test:
	go test -v ./...

clean:
	rm -f rrk

patch:
	@./scripts/bump-version.sh patch
	@if [ -z "$$GITHUB_ACTIONS" ]; then git push --follow-tags; fi

minor:
	@./scripts/bump-version.sh minor
	@if [ -z "$$GITHUB_ACTIONS" ]; then git push --follow-tags; fi

major:
	@./scripts/bump-version.sh major
	@if [ -z "$$GITHUB_ACTIONS" ]; then git push --follow-tags; fi

version:
	@git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
