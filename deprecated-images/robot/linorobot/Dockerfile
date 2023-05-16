ARG BASE_IMAGE="robolaunchio/foxy:plain"

FROM ${BASE_IMAGE}

ARG ROS_DISTRO

SHELL [ "/bin/bash", "-c" ]

ENV ROBOT_WORKSPACE /lino_ws

RUN mkdir -p $ROBOT_WORKSPACE/src && \
    cd $ROBOT_WORKSPACE/src && \
    git clone --single-branch --branch ${ROS_DISTRO} https://github.com/alpylmz/linorobot2.git

RUN source /opt/ros/${ROS_DISTRO}/setup.bash && apt-get update && rosdep init && rosdep update && cd $ROBOT_WORKSPACE && rosdep install --from-path src --ignore-src -y --skip-keys microxrcedds_agent --skip-keys micro_ros_agent

RUN apt-get update && apt-get -y install python3-colcon-common-extensions && apt-get install python3-distutils

RUN source ~/.bashrc && source /opt/ros/${ROS_DISTRO}/setup.bash && cd $ROBOT_WORKSPACE && colcon build --symlink-install

CMD ["bash", "-c", "export DISPLAY=:1 && export LINOROBOT2_BASE=mecanum && source /opt/ros/${ROS_DISTRO}/setup.bash && cd $ROBOT_WORKSPACE && source install/setup.bash && ros2 launch linorobot2_gazebo gazebo.launch.py"]