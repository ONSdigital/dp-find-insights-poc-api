#!/bin/bash -eux

pushd dp-find-insights-poc-api
  make test
  make check-generate
popd
