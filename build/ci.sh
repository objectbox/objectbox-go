#!/usr/bin/env bash
set -euo pipefail

args="$@"

bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-c/main/download.sh) --quiet --sync 0.15.0
export CGO_LDFLAGS="-L$(pwd -P)/lib -Wl,-rpath -Wl,$(pwd -P)/lib"

if [[ "$(uname)" == MINGW* ]]; then
    # copy the dll or ld.exe fails
    cp -vf lib/objectbox.dll /c/TDM-GCC-64/lib/
fi

./build/build.sh $args