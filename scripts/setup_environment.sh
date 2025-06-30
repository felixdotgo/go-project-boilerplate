#!/bin/bash
source $(dirname "$0")/colors.sh

SHELL_NAME=$(echo $SHELL | awk -F '/' '{print $NF}')

# Install Homebrew
if [ -z "$(which brew)" ]; then
  echo "✨ Homebrew is not installed, installing it..."
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

  echo "✨ Setting up Homebrew shell environment..."

  if [ "$SHELL_NAME" = "bash" ]; then
    echo >> $HOME/.bashrc
    echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> $HOME/.bashrc
  elif [ "$SHELL_NAME" = "zsh" ]; then
    echo >> $HOME/.zshrc
    echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> $HOME/.zshrc
  fi

  eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

  echo "✅ Homebrew installation complete."
else
  echo "✅ Homebrew is already installed, skipping installation."
fi

# Install Protobuf
if [ -z "$(which protoc)" ]; then
  echo "✨ Protobuf is not installed, installing it..."
  brew install protobuf@29

  echo "✨ Setting up Protobuf shell environment..."
  if [ "$SHELL_NAME" = "bash" ]; then
    echo >> $HOME/.bashrc
    echo 'export PATH="/home/linuxbrew/.linuxbrew/opt/protobuf@29/bin:$PATH"' >> $HOME/.bashrc
    source $HOME/.bashrc
  elif [ "$SHELL_NAME" = "zsh" ]; then
    echo >> $HOME/.zshrc
    echo 'export PATH="/home/linuxbrew/.linuxbrew/opt/protobuf@29/bin:$PATH"' >> $HOME/.zshrc
    source $HOME/.zshrc
  fi

  echo "✅ Protobuf installation complete."
else
  echo "✅ Protobuf is already installed, skipping installation."
fi

# Install buf tool
if [ -z "$(which buf)" ]; then
  echo "✨ Buf is not installed, installing it..."
  brew install bufbuild/buf/buf
  echo "✅ Buf installation complete."

else
  echo "✅ Buf is already installed, skipping installation."
fi

# Install Go
echo "✨ Installing Go dependencies for generating protobuf code... ($(go version))"

go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
