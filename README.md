# <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/logos/svg/rocket.svg" width="40" height="40" align="top"> Robot Image Pipeline

<div align="center">
  <p align="center">
    <a href="https://github.com/robolaunch/robot-image-pipeline/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/robolaunch/robot-image-pipeline" alt="license">
    </a>
    <a href="https://github.com/robolaunch/robot-image-pipeline/issues">
      <img src="https://img.shields.io/github/issues/robolaunch/robot-image-pipeline" alt="issues">
    </a>
  </p>
</div>

Robot image pipeline produces images for [robolaunch Robot Operator](https://github.com/robolaunch/robot-operator).

## Table of Contents

<!-- - [Overview](#overview) -->
- [Quick Start](#quick-start)
<!-- - [Aims & Roadmap](#aims--roadmap) -->
- [Contributing](#contributing)


<!-- ## Overview

[EDIT THIS: Give more insight about the project. Provide a feature list.]

The aim of this project is to maintain a generic template for robolaunch projects. Members of robolaunch organization can fork this repository and start developing their projects following conventions such as:

- Following a code of conduct
- Having a contributing guide
- Having a style guide
- Applying Apache 2.0 license
- Having a README template
- Having issue & pull request templates
- Using worklows for testing & build -->

## Quick Start

### Running Robot Image Pipeline

See all images from [robolaunch's Docker Hub](https://hub.docker.com/u/robolaunchio).

Clone the artifacts repository.

```bash
git clone https://bitbucket.org/kaesystems/robolaunch-artifacts
```

#### Prerequisites

- `docker-cli`
  - Make sure that you logged in with `docker login`.
- [For VDI Images] Nvidia Driver
  - Check your driver version with `head -n1 </proc/driver/nvidia/version | awk '{print $8}'`
  - You can only build images with your Nvidia driver version.
  - Build VDI `agnostic` if you do not would like to build images with a specific driver version. Installation script (`/etc/vdi/install_driver.sh`) will install appropriate drivers to your system at runtime.


#### Building With Environment Variables

Example image build process:

```bash
export REGISTRY="robolaunchio"
export PUSH_IMAGES="true"
export CONTAINS_VDI="true"
export REBUILD_ROBOT="true"
export ROS_DISTRO="foxy"
export ROBOT_NAME="linorobot"

cd robolaunch-artifacts/robot-operator/images
chmod +x ./build.sh
./build.sh
```

##### `REGISTRY`

Set this environment variable your container registry username. (eg. `robolaunchio`)

##### `PUSH_IMAGES`

- `true`: Every image generated will be pushed to container registry.

- `false`: Images will not be pushed, they will only be built locally.

##### `CONTAINS_VDI`

- `true`: Image(s) contain VDI.

- `false`: Image(s) don't contain VDI.


##### `REBUILD_ROBOT`

- `true`: Builds robolaunch's robot image locally.

- `false`: Tries to pull robolaunch's robot image from Docker Hub.

##### `REBUILD_VDI`

- `true`: Builds VDI image locally.

- `false`: Tries to pull VDI image from Docker Hub.

##### `GPU_AGNOSTIC`

- `true`: Nvidia driver can be installed at runtime instead of in Dockerfile.
  - This option is more convenient if you don't know the driver version in your host system.
  - To install the appropriate Nvidia driver for your system at runtime, run `/etc/vdi/install-driver.sh` as `root` after the container started.

- `false`: Nvidia driver will be installed when building image.


##### `NVIDIA_DRIVER_VERSION`

This variable's value should match your host system and the system which you build your image (to fetch driver modules from host machine). (eg. `470.141.03`)

##### `UBUNTU_DESKTOP`

Select an Ubuntu desktop for VDI. Options are listed under [`robot-operator/images/vdi`](../images/vdi/).

##### `ROS_DISTRO`

Select a ROS distro. Options are listed under [`robot-operator/images/ros`](../images/ros/).

##### `ROBOT_NAME`

Select a robot. Options are listed under [`robot-operator/images/robot`](../images/robot/).

#### Building Interactively

If you do not set any environment variables that configure build script, you can configure your images interactively from the terminal.

```bash
cd robolaunch-artifacts/robot-operator/images
chmod +x ./build.sh
./build.sh
```

#### Adding New Images

If you would like to inject your own images to this pipeline, you should add your Dockerfiles under one of these three directories:

- `robot-operator/images/vdi/<new>`
  - if you want to add an Ubuntu desktop option
- `robot-operator/images/ros/<new>`
  - if you want to add a ROS distro
- `robot-operator/images/robot/<new>`
  - if you want to add a robot
  
<!-- ## Aims & Roadmap

[EDIT THIS: Add roadmap items for the project.]

- Extending the open source conventions
- Enforcing conventional commit messages -->

## Contributing

Please see [this guide](./CONTRIBUTING) if you want to contribute.
