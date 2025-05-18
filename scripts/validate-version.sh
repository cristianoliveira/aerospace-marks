#!/usr/bin/env bash

# Validate Version
# Checks the latest tag in git and compares it with the version in main.go (var VERSION)
#

GIT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.1")
VERSION=$(grep -oP '(?<=var VERSION = ")[^"]*' cmd/main.go)
VERSION_FILE="cmd/main.go"

if [ -z "$VERSION" ]; then
    # IF CI is defined, exit with error
    if [ -n "$CI" ]; then
        echo "ERROR: var VERSION not found in $VERSION_FILE"
        exit 1
    fi

    echo "Version not found in cmd/main.go."
    read -p "Do you want to add the version to cmd/main.go? (y/n) " answer
    if [[ "$answer" == "y" ]]; then

        read -p "[IMPORTANT] Is main release version? Otherwise is nighlty (y/n) " is_main_release
        if [[ "$is_main_release" == "y" ]]; then
            # Version: v0.0.0
            GIT_TAG=$(echo $GIT_TAG | sed 's/^v//')
            GIT_TAG="v$GIT_TAG"
        else
            # Version: v0.0.0-$(YYYYMMDD)-$(GIT_COMMIT)
            GIT_TAG=$(echo $GIT_TAG | sed 's/^v//')
            GIT_TAG="v$GIT_TAG-$(date +%Y%m%d)-$(git rev-parse --short HEAD)"
        fi

        # Remove var VERSION line if it exists  
        sed -i '/var VERSION = /d' $VERSION_FILE
        echo "var VERSION = \"$GIT_TAG\"" >> $VERSION_FILE
        echo "Added version $GIT_TAG to cmd/main.go"
    else
        echo "Version not found and not added. Exiting."
    fi
    exit 1
fi

if [ "$GIT_TAG" != "$VERSION" ]; then
    # IF CI is defined, exit with error
    if [ -n "$CI" ]; then
        echo "ERROR: Version mismatch: GIT_TAG ($GIT_TAG) does not match VERSION ($VERSION)"
        exit 1
    fi

    echo "Version mismatch: GIT_TAG ($GIT_TAG) does not match VERSION ($VERSION)"
    read -p "Do you want to update the version in cmd/main.go to match the GIT_TAG? (y/n) " answer
    if [[ "$answer" == "y" ]]; then
        sed -i "s/var VERSION = \"$VERSION\"/var VERSION = \"$GIT_TAG\"/" $VERSION_FILE
        echo "Updated cmd/main.go to version $GIT_TAG"
    else
        echo "Version mismatch not resolved. Exiting."
    fi

    exit 1
fi

echo "Version match: GIT_TAG ($GIT_TAG) matches VERSION ($VERSION)"
echo "Version check passed."
exit 0
