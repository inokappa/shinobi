#!/usr/bin/env bash
export COGNITO_USER_POOL_ID=$(
  aws --profile=dummy_profile --region=ap-northeast-1 --endpoint=http://192.168.0.100:5000 \
    cognito-idp create-user-pool \
      --pool-name=dummy-user-pool \
      --cli-input-json=file://tests/scripts/dummy-user-pool.json \
      --query=UserPool.Id \
      --output=text)

export COGNITO_CLIENT_ID=$(
  aws --profile=dummy_profile --region=ap-northeast-1 --endpoint=http://192.168.0.100:5000 \
    cognito-idp create-user-pool-client \
      --user-pool-id=${COGNITO_USER_POOL_ID} \
      --client-name=dummy-user-pool-client \
      --cli-input-json=file://tests/scripts/dummy-user-pool-client.json \
      --query=UserPoolClient.ClientId \
      --output=text)