#!/bin/sh

## this script must be called ONLY
## from the root of this project

## this script installs a the
## version of the plugin mentioned
## in the plugin.yaml

version="$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)"
./scripts/install.sh $version