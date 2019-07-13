#!/bin/bash

if [ "${CIRCLE_BRANCH}" = "master" ]; then
  ELS_VER=`git describe --tags  --abbrev=0`
  ELS_NVER="${ELS_VER%.*}.`expr ${ELS_VER##*.} + 1`"
  echo "${ELS_VER} -> ${ELS_NVER}"
  git tag ${ELS_NVER}
else
  ELS_VER=`git describe --tags  --abbrev=0`
  if [[ ${ELS_VER} == *-* ]]; then
    git tag -d ${ELS_VER}
  fi
  ELS_VER=`git describe --tags  --abbrev=0`
  ELS_NVER="${ELS_VER%.*}.`expr ${ELS_VER##*.} + 1`-${CIRCLE_BRANCH}"
  echo "${ELS_VER} -> ${ELS_NVER}"
  git tag ${ELS_NVER}
fi
git push --tags