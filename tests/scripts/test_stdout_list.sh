#!/usr/bin/env bash

source tests/scripts/_setup.sh
gom run shinobi.go -profile=dummy_profile -endpoint=http://192.168.0.100:5000
