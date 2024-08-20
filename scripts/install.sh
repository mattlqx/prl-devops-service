#!/bin/bash
MODE="INSTALL"
INSTALL_SERVICE="true"
while [[ $# -gt 0 ]]; do
  case $1 in
  -i)
    MODE="INSTALL"
    shift # past argument
    ;;
  --install)
    MODE="INSTALL"
    shift # past argument
    ;;
  -u)
    MODE="UNINSTALL"
    shift # past argument
    ;;
  --uninstall)
    MODE="UNINSTALL"
    shift # past argument
    ;;
  -p)
    DESTINATION=$2
    shift # past argument
    shift # past argument
    ;;
  --path)
    DESTINATION=$2
    shift # past argument
    shift # past argument
    ;;
  --no-service)
    INSTALL_SERVICE="false"
    shift # past argument
    ;;
  -v)
    VERSION=$2
    shift # past argument
    shift # past argument
    ;;
  --version)
    VERSION=$2
    shift # past argument
    shift # past argument
    ;;
  *)
    echo "Invalid option $1" >&2
    exit 1
    ;;
  esac
done

if [ -z "$DESTINATION" ]; then
  DESTINATION="/usr/local/bin"
fi

function uninstall() {
  OS=$(uname -s)
  OS=$(echo "$OS" | tr '[:upper:]' '[:lower:]')

  if [ -f "$DESTINATION/prldevops" ]; then
    if [ "$OS" = "darwin" ]; then
      if [ -f "/Library/LaunchDaemons/com.parallels.prl-devops-service.plist" ]; then
        echo "Uninstalling prldevops service"
        echo "Stopping prl-devops-service"
        sudo launchctl unload /Library/LaunchDaemons/com.parallels.prl-devops-service.plist
        sudo rm /Library/LaunchDaemons/com.parallels.prl-devops-service.plist
      fi
    fi

    echo "Removing prldevops from $DESTINATION"
    sudo rm "$DESTINATION/prldevops"
    echo "prldevops has been uninstalled"
  else
    echo "prldevops is not installed in $DESTINATION"
  fi
}

function install() {
  if [ -z "$VERSION" ]; then
    # Get latest version from github
    VERSION=$(curl -s https://api.github.com/repos/Parallels/prl-devops-service/releases/latest | grep -o '"tag_name": "[^"]*"' | cut -d ' ' -f 2 | tr -d '"')
  fi

  if [[ ! $VERSION == *-beta ]]; then
    if [[ ! $VERSION == release-v* ]]; then
      VERSION="release-v$VERSION"
    fi
    SHORT_VERSION="$(echo $VERSION | cut -d '-' -f 2)"
  else
    if [[ ! $VERSION == v* ]]; then
      VERSION="v$VERSION"
    fi
    SHORT_VERSION=$VERSION
  fi

  ARCHITECTURE=$(uname -m)
  if [ "$ARCHITECTURE" = "aarch64" ]; then
    ARCHITECTURE="arm64"
  fi
  if [ "$ARCHITECTURE" = "x86_64" ]; then
    ARCHITECTURE="amd64"
  fi

  OS=$(uname -s)
  OS=$(echo "$OS" | tr '[:upper:]' '[:lower:]')
  echo "Installing prldevops $SHORT_VERSION for $OS-$ARCHITECTURE"

  DOWNLOAD_URL="https://github.com/Parallels/prl-devops-service/releases/download/$VERSION/prldevops--$OS-$ARCHITECTURE.tar.gz"

  echo "Downloading prldevops release from GitHub Releases"
  curl -sL "$DOWNLOAD_URL" -o prldevops.tar.gz

  echo "Extracting prldevops"
  tar -xzf prldevops.tar.gz

  if [ ! -d "$DESTINATION" ]; then
    echo "Creating destination directory: $DESTINATION"
    mkdir -p "$DESTINATION"
  fi

  if [ -f "$DESTINATION/prldevops" ]; then
    echo "Removing existing prldevops"
    sudo rm "$DESTINATION/prldevops"
  fi
  echo "Moving prldevops to $DESTINATION"
  sudo mv prldevops "$DESTINATION"/prldevops
  sudo chmod +x "$DESTINATION"/prldevops

  if [ "$INSTALL_SERVICE" = "true" ]; then
    if [ "$OS" = "darwin" ]; then
      echo "Installing prldevops service"
      if [ -f "/Library/LaunchDaemons/com.parallels.prl-devops-service.plist" ]; then
        echo "Restarting prl-devops-service"
        sudo launchctl unload /Library/LaunchDaemons/com.parallels.prl-devops-service.plist
        sudo launchctl load /Library/LaunchDaemons/com.parallels.prl-devops-service.plist
      fi

      sudo xattr -d com.apple.quarantine "$DESTINATION"/prldevops
    fi
  fi

  echo "Cleaning up"
  rm prldevops.tar.gz
  echo "prldevops $SHORT_VERSION has been installed to $DESTINATION"
}

if [ "$MODE" = "UNINSTALL" ]; then
  uninstall
else
  install
fi
