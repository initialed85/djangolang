#!/bin/bash

set -e

function teardown() {
    ./run.sh env down
}

trap teardown exit

./run.sh env up

./run.sh test-ci
