#!/bin/bash

# バージョンアップスクリプト
# 使用法: ./bump-version.sh [major|minor|patch]

set -e

TYPE=${1:-patch}

# 現在のバージョンタグを取得
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
VERSION=${CURRENT_VERSION#v}

# バージョン番号を分割
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION"

# バージョンタイプに応じて番号を増加
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

# 新しいバージョンを作成
NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"

echo "Bumping version from $CURRENT_VERSION to $NEW_VERSION"

# Gitタグを作成
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"

echo "Created tag $NEW_VERSION"
echo "Run 'git push origin $NEW_VERSION' to trigger release"
