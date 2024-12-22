#!/bin/bash
set -e
SCRIPTS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
. $SCRIPTS_DIR/setenv

export APP_VERSION=$($SCRIPTS_DIR/getver)
export NODE_ENV="production"
rm -rvf "$DIST_DIR"

for PAGE_NAME in res/tmpl/*.prod.html; do
  export PAGE_NAME=$(basename $PAGE_NAME | cut -d. -f1)
  echo "Building $PAGE_NAME"
  bun run internal:tsc && bun run internal:vite build --mode production --emptyOutDir false
done
