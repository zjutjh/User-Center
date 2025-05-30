#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

source hack/util.sh
SCRIPT_PATH=$(util::resolve_path "${BASH_SOURCE[0]}")
SCRIPT_DIR=$(cd "$(dirname "$SCRIPT_PATH")" && pwd -P)
PRODUCT=kantaloupe
PRODUCT_VERSION=v1
OUTPUT_DIR=./api/assets/swagger/v1
DIRS=$(util::get_api_dirs)
PROTO_DIR=./api/
cd "$(dirname "$SCRIPT_PATH")"
cd ../../
rm -rf ${OUTPUT_DIR}/*
mkdir -p ${OUTPUT_DIR}

# Generate and merge into one swagger json file.
protoc \
-I . \
-I ./api/third_party/ \
--openapiv2_out ${OUTPUT_DIR} \
--openapiv2_opt logtostderr=true,allow_delete_body=true \
--openapiv2_opt allow_merge=true \
--openapiv2_opt output_format=json \
--openapiv2_opt merge_file_name="${PRODUCT}.${PRODUCT_VERSION}." \
${PROTO_DIR}/*/*/*.proto ${PROTO_DIR}/*/*.proto

echo "Generated swagger json file: ${OUTPUT_DIR}/${PRODUCT}.${PRODUCT_VERSION}.swagger.json"
cp ./api/third_party/swagger-ui/* ${OUTPUT_DIR}/