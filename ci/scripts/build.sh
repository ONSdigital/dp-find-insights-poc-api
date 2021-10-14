#!/bin/bash -eux

pushd dp-find-insights-poc-api
  make build
  cp build/dp-find-insights-poc-api Dockerfile.concourse ../build
popd
