# Style Guide

Refer [this document](https://github.com/Haufe-Lexware/docker-style-guide/blob/master/Dockerfile.md).

- You should always get base images from an argument named `BASE_IMAGE`.

```dockerfile
ARG BASE_IMAGE=robolaunchio/vdi:focal-agnostic-xfce-withuser
FROM ${BASE_IMAGE} as build
# ...
```