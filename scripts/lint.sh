#!/bin/sh
set -e

# DESIGNED TO BE RUN AS "MAKE LINT"
#
# TODO MAKE SURE RUN IN ROOT DIRECTORY

# https://github.com/golangci/golangci-lint
LINT=./bin/golangci-lint
VERSION=1.27.0

# first time install
if [ ! -f "${LINT}" ]; then 
  echo "Installing ${LINT} to ${VERSION}"
  mkdir -p ./bin
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s "v${VERSION}"
fi

if [ $($LINT --version | grep -c $VERSION) -ne "1" ]; then
  echo "Updating ${LINT} to ${VERSION}"
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s "v${VERSION}"
fi

${LINT} run 

## All of these turned up goodies but not fixing them right now
#
# these are mostly the same
#       --enable=varcheck
#	--enable=deadcode 
#       --enable=unused
#
# useful but will require work
#	--enable=dupl

# // no so useful so far
#	--enable=goconst

# can't tell what it does
#	--enable=scopelint


#	--enable=gochecknoinits
#	--enable=staticcheck \
#	--enable=errchec
#	--enable=unconvert \
#	--enable=gosimple \
