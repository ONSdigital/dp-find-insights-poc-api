#!/bin/bash -eux

pushd dp-find-insights-poc-api
  make check-generate
  make test
popd
