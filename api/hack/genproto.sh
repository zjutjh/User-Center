#!/usr/bin/env bash
#
# Generate all protobuf bindings.
#
set -o pipefail
source hack/util.sh

SCRIPT_PATH=$(util::resolve_path "${BASH_SOURCE[0]}")
SCRIPT_DIR=$(cd "$(dirname "$SCRIPT_PATH")" && pwd -P)
DIRS=$(util::get_api_dirs)
cd "$(dirname "$SCRIPT_PATH")"
cd ../../

echo "Generating proto files..."

CURRENT_VERSION="v1"
GRPC_GATEWAY_TS_OUT=./api/sdk/$CURRENT_VERSION/ts
rm -rf $GRPC_GATEWAY_TS_OUT/*
mkdir -p $GRPC_GATEWAY_TS_OUT

for dir in ${DIRS[@]}
do
  echo "Generating proto for ${dir}"
  # if [ ${dir} != 'v1' ]; then
    for var in `find ./api/${dir} -name "*.proto"`
    do
      clang-format -style="{BasedOnStyle: Google, IndentWidth: 4, ColumnLimit: 0, AlignConsecutiveAssignments: true}" -i "${var}";

      protoc -I . \
      -I ./api/third_party/ \
      --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      --grpc-gateway_out=logtostderr=true,allow_delete_body=true:. --grpc-gateway_opt=paths=source_relative \
      --grpc-gateway-ts_out=$GRPC_GATEWAY_TS_OUT \
      "${var}";
    done
  # fi
done

go mod tidy
echo "Done."
