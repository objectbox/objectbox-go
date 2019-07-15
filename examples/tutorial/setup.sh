#!/bin/bash
set -eu

GO_VER_NEW=1.12.7

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
    if [ -x "$(command -v curl)" ]; then
        curl -L -o go$GO_VER_NEW.linux-$ARCH.tar.gz https://dl.google.com/go/go$GO_VER_NEW.linux-$ARCH.tar.gz
    else
        wget -q https://dl.google.com/go/go$GO_VER_NEW.linux-$ARCH.tar.gz
    fi
    goroot=~/goroot
    echo "Installing go $GO_VER_NEW into $goroot"
    mkdir -p $goroot
    tar -C $goroot -xzf go$GO_VER_NEW.linux-$ARCH.tar.gz --strip 1
    rm go$GO_VER_NEW.linux-$ARCH.tar.gz

    # comment out old GOROOT or GOPATH exports in bashrc if they already exist
    sed -i 's/^export \(GOROOT\|GOPATH\)/#&/' ~/.bashrc
    sed -i 's/^export PATH *=.*\$\(GOROOT\|GOPATH\)/#&/' ~/.bashrc

    # set Go's environment variables in bashrc
    echo "Setting up ~/go as GOPATH"
    mkdir -p ~/go
    echo 'export GOROOT=$HOME/goroot' >> ~/.bashrc
    echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    echo 'export PATH=$GOROOT/bin:$PATH' >> ~/.bashrc
    echo "Note: your ~/.bashrc has been adjusted to make Go $GO_VER_NEW be available globally"
    echo "Please execute \"source ~/.bashrc\" after this script finishes to make Go available for this session or restart your shell"

    # set these variables again for this session
    export GOROOT=$HOME/goroot
    export GOPATH=$HOME/go
    export PATH=$GOROOT/bin:$PATH
}

# setup our working directory
cd ~
mkdir -p projects objectbox

# check if recent version (>=1.12) of Go is installed, otherwise download it
if [ ! -x "$(command -v go)" ]; then
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
cat > update-objectbox.sh <<EOL
#!/bin/bash
set -eu
cd "$(dirname "$0")"
bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-go/master/install.sh)
EOL
chmod +x update-objectbox.sh
./update-objectbox.sh

# create the demo project from examples/tasks
go get -d github.com/objectbox/objectbox-go/examples/tasks
mkdir -p ~/projects/objectbox-go-test/
cd ~/projects/objectbox-go-test/
go mod init objectbox-go-test
cp -r $GOPATH/src/github.com/objectbox/objectbox-go/examples/tasks/* ~/projects/objectbox-go-test/
cd ~/projects/objectbox-go-test
go generate ./...
go build
./objectbox-go-test
