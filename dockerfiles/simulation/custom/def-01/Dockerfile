ARG BASE_IMAGE
FROM ${BASE_IMAGE}

ARG GO1_GYM_PATH

ARG REPOSITORY_PATH="/opt/repositories"

SHELL ["/bin/bash", "-c"]
ARG DEBIAN_FRONTEND=noninteractive

WORKDIR ${REPOSITORY_PATH}

RUN mkdir go1_gym

COPY ${GO1_GYM_PATH}/go1_def_01_demo_code.zip ./go1_gym

RUN set -eux; \
    cd go1_gym; \
    unzip go1_def_01_demo_code.zip -d ./go1_def_01_demo_code; \
    cd go1_def_01_demo_code/go1*/walk-these-ways; \
    pip install -e .

RUN set -eux; \
    pip install opencv-python

WORKDIR /home/robolaunch