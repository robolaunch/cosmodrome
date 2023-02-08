ARG BASE_IMAGE="ros:foxy-ros-base-focal"

FROM golang:1.17 as cli-builder
WORKDIR /go/src
RUN git clone https://bitbucket.org/kaesystems/rlros2 -b main
WORKDIR /go/src/rlros2
RUN make

FROM ${BASE_IMAGE}
ARG BRIDGE_DISTRO_1="foxy"
ARG BRIDGE_DISTRO_2="foxy"

SHELL [ "/bin/bash", "-c" ]
RUN apt update && apt install -y \
    curl \
    pip \
    git \
    wget \
    dnsutils \ 
    ros-${BRIDGE_DISTRO_1}-rosbridge-suite \
    ros-${BRIDGE_DISTRO_2}-rosbridge-suite

# cloud ide & extensions
RUN curl -fsSL https://code-server.dev/install.sh | sh
COPY --from=cli-builder /go/src/rlros2/hack/ros2_cli_override.sh /etc/
COPY --from=cli-builder /go/src/rlros2/rlros2 /usr/bin/

# RUN mkdir -p ~/.local/share/code-server/extensions
# WORKDIR /root/.local/share/code-server/extensions

# ARG CPPTOOLS_VERSION=v1.9.8
# RUN wget https://github.com/microsoft/vscode-cpptools/archive/refs/tags/${CPPTOOLS_VERSION}.tar.gz && tar -x -f ${CPPTOOLS_VERSION}.tar.gz && rm ${CPPTOOLS_VERSION}.tar.gz && mv vscode-cpptools-1.9.8 ms-vscode.cpptools
# # platformio
# ARG PLATFORMIO_VERSION=v2.4.3
# RUN wget https://github.com/platformio/platformio-vscode-ide/archive/refs/tags/${PLATFORMIO_VERSION}.tar.gz && tar -x -f ${PLATFORMIO_VERSION}.tar.gz && rm ${PLATFORMIO_VERSION}.tar.gz
# # robolaunch theme
# RUN git clone https://github.com/gokhangunduz/robolaunch-vscode-theme

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

RUN apt-get dist-upgrade -y
WORKDIR /root
