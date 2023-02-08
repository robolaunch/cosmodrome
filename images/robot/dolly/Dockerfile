ARG BASE_IMAGE="robolaunchio/foxy:plain"

FROM ${BASE_IMAGE}

ARG ROS_DISTRO

SHELL [ "/bin/bash", "-c" ]

ENV ROBOT_WORKSPACE /dolly_ws

RUN mkdir -p $ROBOT_WORKSPACE/src && \
    cd $ROBOT_WORKSPACE/src && \
    git clone --single-branch --branch ${ROS_DISTRO} https://github.com/alpylmz/dolly.git

RUN source /opt/ros/${ROS_DISTRO}/setup.bash && apt-get update && rosdep init && rosdep update && cd $ROBOT_WORKSPACE && rosdep install --from-path src --ignore-src -y --skip-keys microxrcedds_agent --skip-keys micro_ros_agent

RUN apt-get update && apt-get -y install python3-colcon-common-extensions && apt-get install python3-distutils

RUN source ~/.bashrc && source /opt/ros/${ROS_DISTRO}/setup.bash && cd $ROBOT_WORKSPACE && colcon build --symlink-install

CMD ["bash", "-c", "export DISPLAY=:1 && source /opt/ros/${ROS_DISTRO}/setup.bash && source /usr/share/gazebo/setup.bash cd $ROBOT_WORKSPACE && source $ROBOT_WORKSPACE/install/setup.bash && ros2 launch dolly_gazebo dolly.launch.py world:=dolly_city.world"]