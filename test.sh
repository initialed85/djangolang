#!/bin/bash

set -e

./run.sh env up

./run.sh test-ci

./run.sh env down
