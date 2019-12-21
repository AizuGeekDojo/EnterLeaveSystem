#!/bin/bash

git remote set-url origin https://${GITHUB_ACTOR}:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git
echo https://:xxx@github.com/${GITHUB_REPOSITORY}.git
echo ${GITHUB_ACTOR}
BRANCHNAME=${GITHUB_REF##*/}
GITHUB_API_KEY=${GITHUB_TOKEN}

if [ "${BRANCHNAME}" = "master" ]; then
# If branch is master, version up normally.
  ELS_VER=`git describe --tags  --abbrev=0`
  while [[ ${ELS_VER} == *-* ]]; do
    git tag -d ${ELS_VER}
    git push origin :${ELS_VER}
    ELS_VER=`git describe --tags  --abbrev=0`
  done
  ELS_NVER="${ELS_VER%.*}.`expr ${ELS_VER##*.} + 1`"
else
# If branch is not master(pull req), use same version
  ELS_VER=`git describe --tags  --abbrev=0`
  if [[ ${ELS_VER} == *-* ]]; then
echo 10
  # If the version is already used, delete that.
    git tag -d ${ELS_VER}
echo 11
    git push origin :${ELS_VER}
  fi
echo 12
  # Set version
  ELS_VER=`git describe --tags  --abbrev=0`
echo 13
  ELS_NVER="${ELS_VER%.*}.`expr ${ELS_VER##*.} + 1`-${BRANCHNAME}"
fi
echo 14
echo "${ELS_VER} -> ${ELS_NVER}"
echo 15
printf ${ELS_NVER} > tagver
# Set version tag
echo 16
git tag ${ELS_NVER}
echo 17
git push --tags
echo 18
# Release binary
export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED="0"
go env
go run release.go ${BRANCHNAME} ${ELS_NVER}
