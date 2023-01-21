#!/bin/bash

BRANCH=$(git rev-parse --abbrev-ref HEAD)
LATEST_COMMIT=$(git rev-parse HEAD)
LATEST_TAG=$(git describe --tags --always --match 'v*' --abbrev=0)
LATEST_TAG_COMMIT=$(git rev-list -n 1 ${LATEST_TAG})

function get(){
  if [[ ${LATEST_COMMIT} == ${LATEST_TAG} ]]; then
    echo "${BRANCH}.${LATEST_COMMIT}"
    exit 0
  fi

  if [[ ${LATEST_COMMIT} != ${LATEST_TAG_COMMIT} ]]; then
    echo "${BRANCH}.${LATEST_COMMIT}"
    exit 0
  fi

  echo "${LATEST_TAG}"
  exit 0
}

"$@"
