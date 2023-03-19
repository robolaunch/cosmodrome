ARG BASE_IMAGE="ubuntu:focal"

FROM ${BASE_IMAGE}
SHELL [ "/bin/bash", "-c" ]
RUN apt update && apt install -y \
    curl \
    pip \
    git \
    wget \
    dnsutils

# cloud ide & extensions
RUN curl -fsSL https://code-server.dev/install.sh | sh
# COPY --from=cli-builder /go/src/rlros2/hack/ros2_cli_override.sh /etc/
# COPY --from=cli-builder /go/src/rlros2/rlros2 /usr/bin/

# RUN mkdir -p ~/.local/share/code-server/extensions
# WORKDIR /root/.local/share/code-server/extensions

# ARG CPPTOOLS_VERSION=v1.9.8
# RUN wget https://github.com/microsoft/vscode-cpptools/archive/refs/tags/${CPPTOOLS_VERSION}.tar.gz && tar -x -f ${CPPTOOLS_VERSION}.tar.gz && rm ${CPPTOOLS_VERSION}.tar.gz && mv vscode-cpptools-1.9.8 ms-vscode.cpptools
# # platformio
# ARG PLATFORMIO_VERSION=v2.4.3
# RUN wget https://github.com/platformio/platformio-vscode-ide/archive/refs/tags/${PLATFORMIO_VERSION}.tar.gz && tar -x -f ${PLATFORMIO_VERSION}.tar.gz && rm ${PLATFORMIO_VERSION}.tar.gz
# # robolaunch theme
# RUN git clone https://github.com/gokhangunduz/robolaunch-vscode-theme

RUN apt-get dist-upgrade -y
WORKDIR /root
