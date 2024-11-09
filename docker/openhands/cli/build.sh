#!/bin/bash
set -xeu

VERSION="${VERSION:-latest}"
SANDBOX="${SANDBOX:-docker.all-hands.dev/all-hands-ai/runtime:0.13-nikolaik}"
IMAGE="${IMAGE:-openaide/openhands:latest}"

function version() {
	local version

	version=$(git tag --list 'v*' --sort=v:refname --merged|tail -1)
    if [ -z "$version" ]; then
        git rev-parse --short HEAD
    else
        echo "$version"
    fi
}

BUILD_ARGS=(
    --build-arg VERSION="${VERSION}"
    --build-arg SANDBOX="${SANDBOX}"
    --build-arg IMAGE="${IMAGE}"
)

docker build -f Dockerfile "${BUILD_ARGS[@]}" --output "type=local,dest=$PWD/dist" -t openhands-cli-build .
