#!/usr/bin/env bash

set -e
set -u

go tool cover -html=coverage.out

