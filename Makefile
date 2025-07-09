.PHONY: build test clean release-patch release-minor release-major

build:
	go build -o rrk

test:
	go test -v ./...

clean:
	rm -f rrk

patch:
	@./scripts/bump-version.sh patch
	@git push origin $$(git describe --tags --abbrev=0)

minor:
	@./scripts/bump-version.sh minor
	@git push origin $$(git describe --tags --abbrev=0)

major:
	@./scripts/bump-version.sh major
	@git push origin $$(git describe --tags --abbrev=0)

version:
	@git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
