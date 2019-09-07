#!/usr/bin/env bash
set -eu

cLibVersion=0.6.0
os=$(uname)

goVersion=$(go version | cut -d' ' -f 3)
goVersionMajor=$(echo "${goVersion}" | cut -d'.' -f1)
goVersionMinor=$(echo "${goVersion}" | cut -d'.' -f2)
goVersionPatch=$(echo "${goVersion}" | cut -d'.' -f3)

if [[ ! "${goVersionMajor}" == "go1" ]]; then
  echo "Unexpected Go major version ${goVersionMajor}, expecting Go 1.11.4+."
  echo "You can proceed and let us know if you think we should extend the support to your version."
elif [[ "${goVersionMinor}" -lt "11" ]]; then
  echo "Invalid Go version ${goVersion}, at least 1.11.4 required."
  exit 1
elif [[ "${goVersionMinor}" == "11" ]] && [[ "${goVersionPatch}" -lt "4" ]]; then
  echo "Invalid Go version ${goVersion}, at least 1.11.4 required."
  exit 1
fi

# if there's no tty this is probably part of a docker build - therefore we install the c-api explicitly
cLibArgs=
if [[ "$os" != MINGW* ]] && [[ "$os" != CYGWIN* ]]; then
  tty -s || cLibArgs="${cLibArgs} --install"
fi

# install the ObjectBox-C library
mkdir -p objectboxlib && cd objectboxlib
bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-c/master/download.sh) ${cLibArgs} ${cLibVersion}
if [[ "$os" == MINGW* ]] || [[ "$os" == CYGWIN* ]]; then
  localLibDir=$(realpath lib)
  echo "Windows must be able to find the objectbox.dll when executing ObjectBox based programs (even tests)."
  echo "One way to accomplish that is to manually copy the objectbox.dll from ${localLibDir} to C:/Windows/System32/ folder."
  echo "See https://golang.objectbox.io/install#objectbox-library-on-windows for more details."
fi
cd -

if [[ -x "$(command -v ldconfig)" ]]; then
  if ! ldconfig -p | grep -q "libobjectbox."; then
    echo "Installation of the C library failed - ldconfig -p doesn't report libobjectbox. Please try running again."
    exit 1
  fi
fi

# go get flatbuffers (if not using go modules) and objectbox
if [[ ! -f "go.mod" ]]; then
  echo "Your project doesn't seem to be using go modules. Installing FlatBuffers & ObjectBox using go get."
  go get -u github.com/google/flatbuffers/go
  go get -u github.com/objectbox/objectbox-go/objectbox

elif grep -q "module github.com/objectbox/objectbox-go" "go.mod"; then
  echo "Seems like we're running inside the objectbox-go directory itself, skipping ObjectBox Go module installation."

else
  go get -u github.com/objectbox/objectbox-go/objectbox
fi

echo "Installation complete."
echo "You can start using ObjectBox by importing github.com/objectbox/objectbox-go/objectbox in your source code."