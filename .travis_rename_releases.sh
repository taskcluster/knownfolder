#!/bin/bash

VERSION="$(cat "${GOPATH}/src/github.com/taskcluster/knownfolder/main_windows.go" | sed -n 's/.*version = "\(.*\)".*/\1/p')"
export RELEASE_FILE="${TRAVIS_BUILD_DIR}/knownfolder-${VERSION}-${GOOS}-${GOARCH}.exe"
mv "${GOPATH}/bin/${GOOS}_${GOARCH}/knownfolder.exe" "${RELEASE_FILE}"
