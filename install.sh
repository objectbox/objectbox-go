#!/usr/bin/env bash
set -eu

goVersion=$(go version | cut -d' ' -f 3)
goVersionMajor=$(echo "${goVersion}" | cut -d'.' -f1)
goVersionMinor=$(echo "${goVersion}" | cut -d'.' -f2)
goVersionPatch=$(echo "${goVersion}" | cut -d'.' -f3)

if [ ! "${goVersionMajor}" == "go1" ]; then
  echo "Unexpected Go major version ${goVersionMajor}, expecting Go 1.11.4+."
  echo "You can proceed and let us know if you think we should extend the support to your version."
elif [ "${goVersionMinor}" -lt "11" ]; then
  echo "Invalid Go version ${goVersion}, at least 1.11.4 required"
  exit 1
elif [ "${goVersionMinor}" == "11" ] && [ "${goVersionPatch}" -lt "4" ]; then
  echo "Invalid Go version ${goVersion}, at least 1.11.4 required"
  exit 1
fi

# install the ObjectBox-C library
mkdir -p objectboxlib && cd objectboxlib
bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-c/master/download.sh)
cd -

if [ -x "$(command -v ldconfig)" ]; then
  if ! ldconfig -p | grep -q "libobjectbox."; then
    echo "installation of the C library failed - ldconfig -p doesn't report libobjectbox. Please try running again."
    exit 1
  fi
fi

# go get objectbox
if [ -f "go.mod" ] && grep -q "module github.com/objectbox/objectbox-go" "go.mod"; then
  echo "Seems like we're running inside the objectbox-go directory itself, skipping \`go get\` of the objectbox module"
else
  go get -u github.com/objectbox/objectbox-go/objectbox

  if [[ -f "go.mod" ]]; then
    echo "Found go.mod, skipping manual flatbuffers installation as they're already a dependency of the objectbox module"
  fi
fi

echo "Installation complete"