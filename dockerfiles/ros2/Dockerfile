
ARG BASE_IMAGE="ubuntu:20.04"
FROM ${BASE_IMAGE}

SHELL ["/bin/bash", "-c"]
ARG DEBIAN_FRONTEND=noninteractive

ARG ROS_DISTRO_CODE=foxy

RUN apt update && apt install -y \
    gnupg2 \
    lsb-release \
    curl \
    ssh \
    && rm -rf /var/lib/apt/lists/*

RUN curl -sSL https://raw.githubusercontent.com/ros/rosdistro/master/ros.key  -o /usr/share/keyrings/ros-archive-keyring.gpg
RUN echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/ros-archive-keyring.gpg] http://packages.ros.org/ros2/ubuntu $(source /etc/os-release && echo $UBUNTU_CODENAME) main" | tee /etc/apt/sources.list.d/ros2.list > /dev/null

RUN apt update && apt install -y ros-${ROS_DISTRO_CODE}-ros-base

RUN echo "source /opt/ros/${ROS_DISTRO_CODE}/setup.bash" >> ~/.bashrc
RUN apt update && apt install -y \
    python3-rosdep \
    python3-rosinstall \
    python3-rosinstall-generator \
    python3-wstool \
    build-essential \
    python3-colcon-common-extensions \
    && rm -rf /var/lib/apt/lists/*