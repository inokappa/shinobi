#!/usr/bin/env bash

source tests/scripts/_setup.sh
aws --profile=dummy_profile --region=ap-northeast-1 --endpoint=http://192.168.0.100:5000 \
  cognito-idp admin-create-user \
    --user-pool-id=${COGNITO_USER_POOL_ID} \
    --username=testtest1 > /dev/null
gom run shinobi.go -profile=dummy_profile -endpoint=http://192.168.0.100:5000
