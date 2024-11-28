#!/bin/bash

set -xeuo pipefail

##
function init() {
    if [ -d "$1" ]; then
        for script in $(find -L "$1" -type f | sort); do
			if [ -x "$script" ]; then
				echo "Executing $script"
				"$script"
			else
				echo "Skipping $script: not executable"
			fi
        done
    else
        echo "$1 does not exist."
    fi
}

##
init /entrypoint.d

exec "$@"
##