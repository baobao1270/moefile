#!/bin/bash
set -e
SCRIPTS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
. $SCRIPTS_DIR/setenv

rm    -rf $PROJECT_DIR/{dist,bin,.vite}
mkdir -p  $PROJECT_DIR/dist
cp   -rv  $PROJECT_DIR/public/*.go       $PROJECT_DIR/dist
