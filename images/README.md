# Image Pipeline [DEPRECATED]

This document is deprecated. Use `cosmodrome` CLI tool for creating image pipelines.

Prerequisites for building images:

- Docker (with `buildx` plugin)

Clone the repository.

```bash
git clone https://github.com/robolaunch/cosmodrome
```

## Robot Images

### Building & Pushing Images

To build ROS 2 Humble image for [robolaunch Robot Operator](https://github.com/robolaunch/robot-operator), run the following commands:

```bash
# go inside the images
cd cosmodrome/images
# set necessary parameters
export REGISTRY=robolaunchio
export PUSH_IMAGES=true
export CONTAINS_VDI=true
export IMAGE_TYPE=robot
export UBUNTU_DISTRO=jammy
export REBUILD_VDI=true
export GPU_AGNOSTIC=true
export UBUNTU_DESKTOP=xfce
export REBUILD_ROBOT=true
export MULTIPLE_ROS_DISTRO=false
export ROS_DISTRO=humble
# start building
./build.sh
```

## Metrics Exporter Images

```bash
cd cosmodrome/images/metrics
# docker build -t <REGISTRY>/custom-metrics-patcher:<UBUNTU-DISTRO>-v<KUBECTL-VERSION> .
docker build -t robolaunchio/custom-metrics-patcher:focal-v1.24.10 .
```

## Kube Dev Suite Images

This pipeline is under active development.
