#!/bin/bash
set -e
SCRIPTS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
. $SCRIPTS_DIR/setenv

export PAGE_NAME="$1"
shift 1 || true
if [ -z "$PAGE_NAME" ]; then
  echo "Please specify a page to run by adding its name as an argument, or the script will use default page (index)."
fi
exec bun run internal:vite --host 0.0.0.0 "$@"
