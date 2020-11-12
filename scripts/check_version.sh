#!/bin/sh

# Checks plugin.yaml version and also the install version script version

set -e

plugin_yaml=$1
# example: refs/tags/0.0.1
install_version_script=$2
refTag=$3
expected_version=$(echo $refTag | cut -d "/" -f 3)

if [ -z $expected_version ];
then
    echo "git tag value is empty"
    exit 1
fi

if [ ! -e $plugin_yaml ];
then
    echo "File $plugin_yaml does not exist!"
    exit 1
fi

plugin_version=$(yq read $plugin_yaml version)

if [ "$plugin_version" != "$expected_version" ];
then
    echo "Plugin version $plugin_version is not the expected version $expected_version"
    exit 1
fi

if [ ! -e $install_version_script ];
then
    echo "File $install_version_script does not exist!"
    exit 1
fi

set +e

grep -q -F $expected_version $install_version_script
wrong_version=$?

set -e

if [ $wrong_version != 0 ];
then
    echo "$install_version_script does not have the right version"
    exit 1
fi
