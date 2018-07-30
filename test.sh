#!/bin/sh
set -e

docker run --rm --name bcgdv_api --network=airQuality_apiSDN bcgdv/api:latest go test
