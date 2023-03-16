#!/usr/bin/env bash
set -xeuo pipefail

for i in {1..10}; do
  ./test_.sh;
done
