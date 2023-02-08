ARG BASE_IMAGE="robolaunchio/foxy:plain"

FROM ${BASE_IMAGE}

ARG ROS_DISTRO

SHELL [ "/bin/bash", "-c" ]
ENV ROBOT_WORKSPACE /panda_ws

RUN mkdir -p $ROBOT_WORKSPACE/src && \
    cd $ROBOT_WORKSPACE/src && \
    git clone --single-branch --branch ${ROS_DISTRO} https://github.com/UniversalRobots/Universal_Robots_ROS2_Driver

RUN source /opt/ros/${ROS_DISTRO}/setup.bash && apt-get update && rosdep init && rosdep update && cd $ROBOT_WORKSPACE && rosdep install --from-path src --ignore-src -y

RUN apt-get update && apt-get -y install python3-colcon-common-extensions && apt-get install python3-distutils

RUN source ~/.bashrc && source /opt/ros/${ROS_DISTRO}/setup.bash && cd $ROBOT_WORKSPACE && colcon build --cmake-args -DCMAKE_BUILD_TYPE=Release

CMD ["bash", "-c", "export DISPLAY=:1 && source /opt/ros/${ROS_DISTRO}/setup.bash && cd $ROBOT_WORKSPACE && source install/setup.bash && ros2 launch ur_bringup ur_control.launch.py ur_type:=ur3 robot_ip:=127.0.0.1 use_fake_hardware:=true launch_rviz:=true"]