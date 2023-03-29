#!/bin/bash

set eux;

if ! command -v nvidia-smi &> /dev/null
then
    echo "nvidia-smi command not found!";
    exit 1;
fi

if [[ -z "${GPU_LATENCY}" ]]; then
    echo "Environment GPU_LATENCY should be set.";
    exit 1;
fi

while [ true ]
do
    nvidia-smi | grep "%" | awk '{print $13}';
    sleep "$GPU_LATENCY";
    # kubectl patch metricsexporter metricsexporter-sample --type=merge -n test --patch '{"spec":{"foo":"asd"}}'
done
