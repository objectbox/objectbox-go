#!/usr/bin/env bash
set -euo pipefail

# macOS does not have realpath and readlink does not have -f option, so do this instead:
script_dir=$( cd "$(dirname "$0")" ; pwd -P )
cd "${script_dir}/.." # move to project root dir

args="$@"

bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-c/main/download.sh) --quiet --sync 4.1.0
export CGO_LDFLAGS="-L$(pwd -P)/lib -Wl,-rpath -Wl,$(pwd -P)/lib"

if [[ "$(uname)" == MINGW* ]]; then
    # copy the dll or ld.exe fails
    cp -vf lib/objectbox.dll /c/TDM-GCC-64/lib/
fi

./build/build.sh $args