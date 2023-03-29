#!/bin/bash

set eux;

if [[ -z "${NETWORK_INTERFACES}" ]]; then
    echo "Environment NETWORK_INTERFACES should be set.";
    exit 1;
fi

ifstat -nb -i "$NETWORK_INTERFACES" 2 5;