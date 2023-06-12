ARG BASE_IMAGE
FROM ${BASE_IMAGE} as build
ARG PHP_VERSION
RUN apt-get update && apt-get install -qy php=${PHP_VERSION} libapache2-mod-php
WORKDIR /home/robolaunch
USER root

