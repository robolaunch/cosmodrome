ARG BASE_IMAGE=nvidia/opengl:1.2-glvnd-runtime-ubuntu20.04

#
# STAGE 1: SERVER
#
FROM ubuntu:focal as server
WORKDIR /root

RUN set -eux; apt-get update; apt-get install -y curl; \
    curl -L https://go.dev/dl/go1.18.5.linux-amd64.tar.gz | tar zx; \
    cp -r ./go /usr/local/

WORKDIR /src

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ UTC

#
# install dependencies
RUN set -eux; apt-get update; \
    apt-get install -y --no-install-recommends build-essential git cmake make libx11-dev libxrandr-dev libxtst-dev \
    libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly; \
    #
    # install libclipboard
    set -eux; \
    cd /tmp; \
    git clone --depth=1 https://github.com/jtanx/libclipboard; \
    cd libclipboard; \
    cmake .; \
    make -j4; \
    make install; \
    rm -rf /tmp/libclipboard; \
    #
    # clean up
    apt-get clean -y; \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/*

#
# build server
# COPY server/ .
RUN git clone https://bitbucket.org/kaesystems/vdi /root/vdi/; \
    cp -a /root/vdi/server/* .

RUN /usr/local/go/bin/go get -v -t -d . && /usr/local/go/bin/go build -o bin/neko cmd/neko/main.go


#
# STAGE 2: CLIENT
#
FROM node:14-bullseye-slim as client
WORKDIR /src

#
# install dependencies

RUN set -eux; apt-get update; \
    apt-get install -y --no-install-recommends git

# COPY client/package*.json ./

ENV GIT_SSL_NO_VERIFY="1"

RUN git clone https://bitbucket.org/kaesystems/vdi /root/vdi/; \
    cp /root/vdi/client/package*.json ./

RUN npm install
#
# build client
# COPY client/ .
RUN cp -a /root/vdi/client/. .

RUN npm run build

#
# STAGE 3: RUNTIME
#
#FROM debian:bullseye-slim
#######################################START

FROM $BASE_IMAGE

LABEL maintainer "https://github.com/ehfd"

ARG NVIDIA_VISIBLE_DEVICES=all
ARG DEBIAN_FRONTEND=noninteractive
ENV NVIDIA_DRIVER_CAPABILITIES all

ENV TZ UTC
ENV PASSWD mypasswd
ENV CDEPTH 24
ENV VIDEO_PORT DFP

#ENV NEKO_PASSWORD secret
#ENV NEKO_PASSWORD_ADMIN secret

#######################################END
#
# avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

RUN set -eux; \
    #
    # add non-free repo for intel drivers
    #echo deb http://deb.debian.org/debian bullseye main contrib non-free > /etc/apt/sources.list; \
    #echo deb http://deb.debian.org/debian-security/ bullseye-security main contrib non-free >> /etc/apt/sources.list; \
    #echo deb http://deb.debian.org/debian bullseye-updates main contrib non-free >> /etc/apt/sources.list; \
    apt-get update; \
    #
    # install dependencies
    apt-get install -y --no-install-recommends nano kmod curl wget ca-certificates supervisor sudo; \
    apt-get install -y --no-install-recommends openssh-server pulseaudio dbus-x11; \
    apt-get install -y --no-install-recommends libcairo2 libxcb1 libxrandr2 libxv1 libopus0; \ 
    # libvpx6;
    #
    # gst + vaapi plugin
    apt-get install -y --no-install-recommends libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev \
                    gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-pulseaudio \
                    gstreamer1.0-vaapi ;\
    #
    # fonts
    apt-get install -y --no-install-recommends fonts-takao-mincho fonts-wqy-zenhei; \
    #
    mkdir -p /var/www /var/log/neko; \
    # clean up
    apt-get clean -y; \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/*

#
# copy config files
COPY ./supervisord.conf /etc/vdi/supervisord.conf
COPY ./install-driver.sh /etc/vdi/install-driver.sh
COPY ./dbus /usr/bin/dbus
COPY ./default.pa /etc/pulse/default.pa
COPY ./add-render-group.sh /usr/bin/add-render-group.sh

RUN chmod 755 /etc/vdi/install-driver.sh

ENV NEKO_PASSWORD=robolaunch
ENV NEKO_PASSWORD_ADMIN=admin
ENV RENDER_GID=

ARG APP_CACHEBUST=1
COPY --from=client /src/dist/ /var/www/
COPY --from=server /src/bin/neko /usr/bin/neko

HEALTHCHECK --interval=10s --timeout=5s --retries=8 \
    CMD wget -O - http://localhost:${NEKO_BIND#*:}/health || exit 1

COPY ./*.so /usr/lib/x86_64-linux-gnu/gstreamer-1.0/

ARG NVIDIA_DRIVER_VERSION="470.141.03"
ENV NVIDIA_DRIVER_VERSION=${NVIDIA_DRIVER_VERSION}
ARG GPU_AGNOSTIC="true"
RUN GPU_AGNOSTIC=${GPU_AGNOSTIC} /etc/vdi/install-driver.sh

CMD ["/usr/bin/supervisord", "-c", "/etc/vdi/supervisord.conf"]