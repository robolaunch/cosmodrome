ARG BASE_IMAGE="ros:foxy-ros-base-focal"

# FROM golang:1.17 as cli-builder
# WORKDIR /go/src
# RUN git clone https://bitbucket.org/kaesystems/rlros2 -b main
# WORKDIR /go/src/rlros2
# RUN make

FROM ${BASE_IMAGE}
ARG BRIDGE_DISTRO_1="foxy"
ARG BRIDGE_DISTRO_2="foxy"

SHELL [ "/bin/bash", "-c" ]
RUN apt-get update && apt install -y \
    curl \
    pip \
    git \
    wget \
    dnsutils \ 
    ros-${BRIDGE_DISTRO_1}-rosbridge-suite \
    ros-${BRIDGE_DISTRO_2}-rosbridge-suite

# tracker & robolaunch.py helpers
ARG TRACKER_PATH=/var/lib/ros-tracker
ARG TRACKER_REPO_URL=https://bitbucket.org/kaesystems/tracker/src/main/
ARG TRACKER_BRANCH=main

ARG ROBOLAUNCH_HELPERS_PATH=/var/lib/robolaunch-helpers
ARG ROBOLAUNCH_HELPERS_REPO_URL=https://bitbucket.org/kaesystems/robolaunch-cli/src/master/
ARG ROBOLAUNCH_HELPERS_BRANCH=master

ARG CACHEBUST=1
RUN git clone ${TRACKER_REPO_URL} ${TRACKER_PATH} -b ${TRACKER_BRANCH}
WORKDIR ${TRACKER_PATH}
RUN pip install --no-cache-dir --upgrade -r requirements.txt

RUN git clone ${ROBOLAUNCH_HELPERS_REPO_URL} ${ROBOLAUNCH_HELPERS_PATH} -b ${ROBOLAUNCH_HELPERS_BRANCH}
WORKDIR ${ROBOLAUNCH_HELPERS_PATH}
RUN pip install --no-cache-dir --upgrade -r requirements.txt

WORKDIR /
RUN echo "Built and packaged: $(date -d @$(date -u +%s))" > image.info

WORKDIR /root
