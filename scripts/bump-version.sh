#!/bin/bash

set -e

TYPE=${1:-patch}

CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
VERSION=${CURRENT_VERSION#v}

IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION"

case $TYPE in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo "Usage: $0 [major|minor|patch]"
        exit 1
        ;;
esac

NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"

echo "Bumping version from $CURRENT_VERSION to $NEW_VERSION"

git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"

echo "Created tag $NEW_VERSION"
echo "Run 'git push origin $NEW_VERSION' to trigger release"
