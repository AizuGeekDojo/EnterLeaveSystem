#!/bin/bash

ssh -T git@github.com -o StrictHostKeyChecking=no
if [ "${CIRCLE_BRANCH}" = "master" ]; then
  ELS_VER=`git describe --tags  --abbrev=0`
  ELS_NVER="${ELS_VER%.*}.`expr ${ELS_VER##*.} + 1`"
else
  ELS_VER=`git describe --tags  --abbrev=0`
  if [[ ${ELS_VER} == *-* ]]; then
    git tag -d ${ELS_VER}
  fi
  ELS_VER=`git describe --tags  --abbrev=0`
  ELS_NVER="${ELS_VER%.*}.`expr ${ELS_VER##*.} + 1`-${CIRCLE_BRANCH}"
fi
echo "${ELS_VER} -> ${ELS_NVER}"
git tag ${ELS_NVER}
git push --tags
go run release.go ${CIRCLE_BRANCH} ${ELS_NVER}