ARG BASE_IMAGE
FROM ${BASE_IMAGE} as build
ARG CPP_VERSION
USER root
RUN apt-get update && apt-get install -qy gcc-9=${CPP_VERSION} git
WORKDIR /home/robolaunch
USER root
