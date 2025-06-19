#!/bin/bash

if "$(git diff HEAD^ HEAD --name-only | grep -v '^nix/')"; then
    echo "true"
else
    echo "false"
fi
