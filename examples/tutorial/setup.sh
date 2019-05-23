#!/bin/bash

GO_VER_NEW=1.12.5

function fetch_go {
    ARCH=$(uname -m)
    if [[ "$ARCH" == "x86_64" ]]; then
        ARCH=amd64
    elif [[ "$ARCH" == arm* ]]; then
        ARCH=armv6l
    else
        echo "Unsupported architecture: $ARCH"
        if [[ "$GO_DIR_MISSING" == "1" ]]; then
            rm -rf ~/go
        fi
    fi

    # download Go
    echo "Fetching go $GO_VER_NEW..."
    wget -q https://dl.google.com/go/go$GO_VER_NEW.linux-$ARCH.tar.gz
    tar xzf go$GO_VER_NEW.linux-$ARCH.tar.gz
    rm go$GO_VER_NEW.linux-$ARCH.tar.gz

    # comment out old GOROOT or GOPATH exports in bashrc if they already exist
    sed -i 's/^export \(GOROOT\|GOPATH\)/#&/' ~/.bashrc
    sed -i 's/^export PATH *=.*\$\(GOROOT\|GOPATH\)/#&/' ~/.bashrc

    # set Go's environment variables in bashrc
    mkdir -p ~/go/libs
    echo 'export GOROOT=$HOME/go/go' >> ~/.bashrc
    echo 'export GOPATH=$HOME/go/libs' >> ~/.bashrc
    echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.bashrc
    echo "Note: your ~/.bashrc has been adjusted to make Go $GO_VER_NEW be available globally"
    echo "Please execute \"source ~/.bashrc\" after this script finishes to make Go available for this session or restart your shell"

    # set these variables again for this session
    export GOROOT=$HOME/go/go
    export GOPATH=$HOME/go/libs
    export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
}

# setup our working directory
cd ~
if [ ! -d "go" ]; then
    GO_DIR_MISSING=1
fi
mkdir -p go
cd go
mkdir -p projects objectbox

# check if recent version (>=1.12) of Go is installed, otherwise download it
command -v go > /dev/null
if (( "$?" != "0" )); then
    fetch_go
else
    GO_VER=$(go version | cut -d' ' -f3)
    GO_VER_MAJOR=$(echo $GO_VER | cut -d'.' -f1 | cut -d'o' -f2)
    GO_VER_MINOR=$(echo $GO_VER | cut -d'.' -f2)
    if (( "$GO_VER_MAJOR" >= "1" && "$GO_VER_MINOR" >= "12" )); then
        echo "Note: using installed Go with version $GO_VER"
    else
        echo "Warning: an old version of Go ($GO_VER) is installed on your system, ObjectBox needs >= 1.12"
        read -n 1 -s -r -p "Would you like to download Go $GO_VER_NEW? (y/n) " inp
        if [[ "$inp" == "y" ]]; then
            fetch_go
        else
            echo "Exiting, as ObjectBox cannot be used with Go $GO_VER"
            if [[ "$GO_DIR_MISSING" == "1" ]]; then
                rm -rf ~/go
            fi
        fi
    fi
fi

# get the ObjectBox binary library
cd objectbox
IFS=
read -d '' updateobjectbox <<"EOF"
#!/bin/bash
curl -s https://raw.githubusercontent.com/objectbox/objectbox-c/master/download.sh 2> /dev/null > download.sh
chmod +x download.sh
./download.sh --quiet 2>&1 > /dev/null
rm download.sh
if [ -f /usr/local/lib/libobjectbox.so ]; then
    echo "The ObjectBox library already exists in /usr/local/lib/libobjectbox.so."
    read -p "Would you like to replace it with the possibly newer version that was just downloaded? (Y/n) " inp
    if [[ "$inp" == "y" || "$inp" == "" ]]; then
        sudo mv lib/libobjectbox.so /usr/local/lib/
    fi
else
    echo "The ObjectBox library needs to be globally available for ObjectBox Go to work correctly."
    read -p "Would you like to install the ObjectBox library system-wide (this might prompt for your root password)? (Y/n) " inp
    if [[ "$inp" == "y" || "$inp" == "" ]]; then
        sudo mv lib/libobjectbox.so /usr/local/lib/
    else
        echo "The library was not installed system-wide, but was moved to the directory ~/go/objectbox instead."
        echo "Later calls to ObjectBox Go might fail, unless you adjust your CGO_LDFLAGS and LD_LIBRARY_PATH environment variables."
        mv lib/libobjectbox.so .
    fi
fi
rm -r download lib
EOF
printf $updateobjectbox > update-objectbox.sh
chmod +x update-objectbox.sh
./update-objectbox.sh

# get the ObjectBox Go library
go get -v github.com/objectbox/objectbox-go/...
go get github.com/google/flatbuffers/go
go install github.com/objectbox/objectbox-go/cmd/objectbox-gogen/

# create the demo project
cd ../projects
mkdir objectbox-go-test
cd objectbox-go-test
wget -q https://raw.githubusercontent.com/objectbox/objectbox-go/dev/examples/tutorial/main.go
objectbox-gogen -source main.go
go build
./objectbox-go-test
