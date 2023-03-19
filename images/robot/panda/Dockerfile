ARG BASE_IMAGE="robolaunchio/foxy:plain"

FROM ${BASE_IMAGE}

ARG ROS_DISTRO

SHELL [ "/bin/bash", "-c" ]
ENV ROBOT_WORKSPACE /panda_ws

RUN mkdir -p $ROBOT_WORKSPACE/src && \
    cd $ROBOT_WORKSPACE/src && \
    git clone --single-branch --branch develop https://github.com/frankaemika/franka_ros2

RUN source /opt/ros/${ROS_DISTRO}/setup.bash && apt-get update && rosdep init && rosdep update && cd $ROBOT_WORKSPACE && rosdep install --from-path src --ignore-src -y

RUN apt-get update && apt-get -y install python3-colcon-common-extensions && apt-get install python3-distutils

# current version of franka_ros2 has this dependency lacked
RUN apt-get update && apt-get -y install ros-${ROS_DISTRO}-control-toolbox

RUN source ~/.bashrc && source /opt/ros/${ROS_DISTRO}/setup.bash && cd $ROBOT_WORKSPACE && colcon build --symlink-install

CMD ["bash", "-c", "export DISPLAY=:1 && source /opt/ros/${ROS_DISTRO}/setup.bash && cd $ROBOT_WORKSPACE && source install/setup.bash && ros2 launch franka_moveit_config moveit.launch.py robot_ip:=dont-care use_fake_hardware:=true"]