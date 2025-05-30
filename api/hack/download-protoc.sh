#!/bin/bash

set -ex

FILE_NAME=protoc-29.3-linux-x86_64.zip
GOOS=$(go env GOOS)
if [[ "$GOOS" == "darwin" ]]; then
	 	FILE_NAME=protoc-29.3-osx-x86_64.zip
fi
GOARCH=$(go env GOARCH)
if [[ "$GOARCH" == "arm64" ]]; then
	 	FILE_NAME=protoc-29.3-osx-aarch_64.zip
fi

STORCLI_URL="https://github.com/protocolbuffers/protobuf/releases/download/v29.3/$FILE_NAME"

curl -L -f "$STORCLI_URL" > $FILE_NAME

mkdir protobuf
unzip $FILE_NAME -d protobuf

cd protobuf
sudo cp ./bin/* /usr/local/bin/
sudo cp -r ./include/* /usr/local/include/

rm ../$FILE_NAME
rm -rf ../protobuf


