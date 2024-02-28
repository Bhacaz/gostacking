#!/bin/zsh

version=$(<VERSION)

sed -i '' "s/Version: \".*\"/Version: \"$version\"/" cmd/root.go

git add .
git commit -m "Release $version"
git tag -a "$version" -m "Release $version"
git push --tags

goreleaser release --clean
