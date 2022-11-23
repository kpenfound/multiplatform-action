#!/bin/bash
set -eux

USER=kpenfound
REPO=typescript-multiplatform-action
GITHUB_API=https://api.github.com/repos/${USER}/${REPO}/releases/latest

ARCH=$(dpkg --print-architecture) 
echo "Querying GitHub for the latest $ARCH release" 
LATEST=$(curl -L -s -H 'Accept: application/json' $GITHUB_API) 
LATEST_TAG=$(echo $LATEST | jq -r '.tag_name') 
echo "Found version $LATEST_TAG"
LATEST_URL=$(echo $LATEST | jq -r ".assets[] | select(.name | contains(\"linux_$ARCH\")) | .url") 
echo "Downloading curl $LATEST_TAG from $USER"
curl -vLJ -o action -H 'Accept: application/octet-stream' $LATEST_URL
chmod +x action
./action