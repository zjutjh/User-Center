#!/usr/bin/env bash
#
# Clean all things by auto generate.
#
set -ex

cleanFunc(){
  for var in `find $1 -name "*.pb.go"`
    do
      rm -rf "$var"
    done
  for var in `find $1 -name "*.pb.gw.go"`
    do
      rm -rf "$var"
    done
}

cleanFunc .
rm -rf sdk
rm -rf assets