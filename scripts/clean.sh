#!/bin/sh
set -ex

# normally I have done go clean ./...
# but under golang modules it seems to work
# differently.  Use
# go clean .

go clean -i -cache -testcache -modcache .

# remove coverprofile from scripts/coverage
rm -f coverprofile.out

rm -rf ./bin
