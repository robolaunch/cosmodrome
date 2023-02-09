#!/bin/bash

function build_vdi() {
    
    ask_gpu_agnostic

    if [ "$GPU_AGNOSTIC" != "true" ]; then
        ask_driver
    else
        NVIDIA_DRIVER_VERSION="agnostic"
    fi

    ask_ubuntu_desktop

    DRIVER_IMAGE=$REGISTRY/driver:$UBUNTU_DISTRO-$NVIDIA_DRIVER_VERSION
    VDI_IMAGE=$REGISTRY/vdi:$UBUNTU_DISTRO-$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP

    docker_build_driver
    docker_build_vdi

}

function build_ros() {

    if [ "$CONTAINS_VDI" == "true" ]; then

        ask_gpu_agnostic
        if [ "$GPU_AGNOSTIC" != "true" ]; then
            ask_driver
        else
            NVIDIA_DRIVER_VERSION="agnostic"
        fi
        ask_ubuntu_desktop

        ask_multiple_ros_distro
        if [ "$MULTIPLE_ROS_DISTRO" == "true" ]; then
            ask_ros_distro

            ROS_BASE_IMAGE=$REGISTRY/vdi:$UBUNTU_DISTRO-$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
            ROS_IMAGE=$REGISTRY/$ROS_DISTRO_1-$ROS_DISTRO_2:$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
        else
            ask_ros_distro

            ROS_BASE_IMAGE=$REGISTRY/vdi:$UBUNTU_DISTRO-$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
            ROS_IMAGE=$REGISTRY/$ROS_DISTRO:$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
        fi
       
    else

        ask_multiple_ros_distro
        if [ "$MULTIPLE_ROS_DISTRO" == "true" ]; then
            ask_ros_distro
        
            ROS_BASE_IMAGE="ubuntu:$UBUNTU_DISTRO"
            ROS_IMAGE=$REGISTRY/$ROS_DISTRO_1-$ROS_DISTRO_2:plain
        else
            ask_ros_distro
        
            ROS_BASE_IMAGE="ubuntu:$UBUNTU_DISTRO"
            ROS_IMAGE=$REGISTRY/$ROS_DISTRO:plain
        fi
        
    fi

    docker_build_ros

}

function build_robolaunch_robot() {

    if [ "$CONTAINS_VDI" == "true" ]; then

        if [ "$MULTIPLE_ROS_DISTRO" == "true" ]; then

            ROBOLAUNCH_ROBOT_IMAGE=$REGISTRY/robot:base-$ROS_DISTRO_1-$ROS_DISTRO_2-$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
            ROBOLAUNCH_ROBOT_BASE_IMAGE=$REGISTRY/$ROS_DISTRO_1-$ROS_DISTRO_2:$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
            
        else

            ROBOLAUNCH_ROBOT_IMAGE=$REGISTRY/robot:base-$ROS_DISTRO-$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
            ROBOLAUNCH_ROBOT_BASE_IMAGE=$REGISTRY/$ROS_DISTRO:$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
            
        fi
        
    else

        if [ "$MULTIPLE_ROS_DISTRO" == "true" ]; then
            
            ROBOLAUNCH_ROBOT_IMAGE=$REGISTRY/robot:base-$ROS_DISTRO_1-$ROS_DISTRO_2-plain
            ROBOLAUNCH_ROBOT_BASE_IMAGE=$REGISTRY/$ROS_DISTRO_1-$ROS_DISTRO_2:plain

        else

            ROBOLAUNCH_ROBOT_IMAGE=$REGISTRY/robot:base-$ROS_DISTRO-plain
            ROBOLAUNCH_ROBOT_BASE_IMAGE=$REGISTRY/$ROS_DISTRO:plain
            
        fi
    
    fi

    docker_build_robolaunch_robot

}

function build_robot() {
    
    if [ "$CONTAINS_VDI" == "true" ]; then
        
        ask_gpu_agnostic
        if [ "$GPU_AGNOSTIC" != "true" ]; then
            ask_driver
        else
            NVIDIA_DRIVER_VERSION="agnostic"
        fi
        ask_ubuntu_desktop
        ask_ros_distro
        ask_robot_name

        if [ "$MULTIPLE_ROS_DISTRO" == "true" ]; then
            
            ask_robot_distro
            ROBOT_IMAGE=$REGISTRY/$ROBOT_NAME:$ROS_DISTRO_1-$ROS_DISTRO_2-$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
            ROBOT_BASE_IMAGE=$REGISTRY/robot:base-$ROS_DISTRO_1-$ROS_DISTRO_2-$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP

        else

            ROBOT_IMAGE=$REGISTRY/$ROBOT_NAME:$ROS_DISTRO-$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
            ROBOT_BASE_IMAGE=$REGISTRY/robot:base-$ROS_DISTRO-$NVIDIA_DRIVER_VERSION-$UBUNTU_DESKTOP
            
        fi
        
    else

        ask_ros_distro
        ask_robot_name

        if [ "$MULTIPLE_ROS_DISTRO" == "true" ]; then
            
            ask_robot_distro
            ROBOT_IMAGE=$REGISTRY/$ROBOT_NAME:$ROS_DISTRO_1-$ROS_DISTRO_2
            ROBOT_BASE_IMAGE=$REGISTRY/robot:base-$ROS_DISTRO_1-$ROS_DISTRO_2-plain

        else

            ROBOT_IMAGE=$REGISTRY/$ROBOT_NAME:$ROS_DISTRO
            ROBOT_BASE_IMAGE=$REGISTRY/robot:base-$ROS_DISTRO-plain
            
        fi
        
    fi

    docker_build_robot
    
}

function ask_ubuntu_distro(){
    
    if [[ -z "${UBUNTU_DISTRO}" ]]; then
        
        # ask nvidia driver version
        echo -n "Select Ubuntu distro. (eg. focal, jammy): "
        read UBUNTU_DISTRO

    fi

    if ! [[ $UBUNTU_DISTRO =~ ^(focal|jammy)$ ]]; then 
            
        echo "Ubuntu version should be focal or jammy."
        exit 1;

    fi

    if [ "$UBUNTU_DISTRO" == "focal" ]; then

        OPENGL_IMAGE="nvidia/opengl:1.2-glvnd-runtime-ubuntu20.04"

    fi

     if [ "$UBUNTU_DISTRO" == "jammy" ]; then

        OPENGL_IMAGE="robolaunchio/opengl:1.2-runtime-ubuntu22.04"

    fi


    
}

function ask_driver(){
    
    if [[ -z "${NVIDIA_DRIVER_VERSION}" ]]; then
        
        # ask nvidia driver version
        echo -n "Select Nvidia driver version. (eg. 470.141.03): "
        read NVIDIA_DRIVER_VERSION

    fi

    if ! [[ $NVIDIA_DRIVER_VERSION =~ ^([[:xdigit:]]{3}.[[:xdigit:]]{3}.[[:xdigit:]]{2}) ]]; then 
            
        echo "Driver version is in wrong format. (eg. 444.555.66)"
        exit 1;

    fi
    
}

function ask_gpu_agnostic() {

    if [[ -z "${GPU_AGNOSTIC}" ]]; then
        
        # ask nvidia driver version
        echo -n "Is image GPU-agnostic? (true/false): "
        read GPU_AGNOSTIC

    fi

    if ! [[ "$GPU_AGNOSTIC" == "true" || "$GPU_AGNOSTIC" == "false" ]]; then

        echo "Wrong input for \"GPU_AGNOSTIC\". (true/false)"
        exit 1;

    fi

}

function ask_ubuntu_desktop(){
    
    if [[ -z "${UBUNTU_DESKTOP}" ]]; then
       
        # ask ubuntu desktop choice
        cd vdi
        echo -n "Select Ubuntu desktop. ("
        for i in *; do if [ "$i" != "base" ]; then echo -n "$i/"; fi done
        echo -n "): "
        read UBUNTU_DESKTOP
        cd ..

    fi

    if ! [[ -d "./vdi/$UBUNTU_DESKTOP" ]]; then

        echo "Ubuntu desktop \"$UBUNTU_DESKTOP\" does not exist."
        exit 1;

    fi

}

function ask_multiple_ros_distro(){
   
    if [[ -z "${MULTIPLE_ROS_DISTRO}" ]]; then
        
        # ask nvidia driver version
        echo -n "Does image contain multiple ROS distro? (true/false): "
        read MULTIPLE_ROS_DISTRO

    fi

    if ! [[ "$MULTIPLE_ROS_DISTRO" == "true" || "$MULTIPLE_ROS_DISTRO" == "false" ]]; then

        echo "Wrong input for \"MULTIPLE_ROS_DISTRO\". (true/false)"
        exit 1;

    fi

}

function ask_ros_distro(){

    if [ "$MULTIPLE_ROS_DISTRO" == "true" ]; then
        
        if [[ -z "${ROS_DISTRO_1}" ]]; then
            # ask ros distro
            cd ros
            echo -n "Select first ROS distro. ("
            for i in *; do echo -n "$i/"; done
            echo -n "): "
            read ROS_DISTRO_1
            cd ..

        fi

        if ! [[ -d "./ros/$ROS_DISTRO_1" ]]; then

            echo "Distro \"$ROS_DISTRO_1\" does not exist."
            exit 1;
        
        fi


        if [[ -z "${ROS_DISTRO_2}" ]]; then
            # ask ros distro
            cd ros
            echo -n "Select second ROS distro. ("
            for i in *; do echo -n "$i/"; done
            echo -n "): "
            read ROS_DISTRO_2
            cd ..

        fi

        if ! [[ -d "./ros/$ROS_DISTRO_2" ]]; then

            echo "Distro \"$ROS_DISTRO_2\" does not exist."
            exit 1;
        
        fi

    else

        if [[ -z "${ROS_DISTRO}" ]]; then
            # ask ros distro
            cd ros
            echo -n "Select ROS distro. ("
            for i in *; do echo -n "$i/"; done
            echo -n "): "
            read ROS_DISTRO
            cd ..

        fi

        if ! [[ -d "./ros/$ROS_DISTRO" ]]; then

            echo "Distro \"$ROS_DISTRO\" does not exist."
            exit 1;
        
        fi

    fi   
    
}

function ask_robot_distro(){


    if [[ -z "${ROBOT_DISTRO}" ]]; then
        # ask robot distro
        cd ros
        echo -n "Select robot's ROS distro. ("
        for i in *; do echo -n "$i/"; done
        echo -n "): "
        read ROBOT_DISTRO
        cd ..

    fi

    if ! [[ -d "./ros/$ROBOT_DISTRO" ]]; then

        echo "Distro \"$ROBOT_DISTRO\" does not exist."
        exit 1;
    
    fi

    
}

function ask_robot_name(){
    
    if [[ -z "${ROBOT_NAME}" ]]; then
        
        # ask robot name
        cd robot
        echo -n "Select robot. ("
        for i in *; do if [ "$i" != "Dockerfile" ]; then echo -n "$i/"; fi done
        echo -n "): "
        read ROBOT_NAME
        cd ..

    fi

    if ! [[ -d "./robot/$ROBOT_NAME" ]];  then

        echo "Robot \"$ROBOT_NAME\" does not exist."
        exit 1;

    fi

}

function ask_registry(){

    # ask if images will be pushed to registry
    if [[ -z "${REGISTRY}" ]]; then

        echo -n "Container registry (Docker Hub) username: (eg. robolaunchio): "
        read REGISTRY

    fi

}

function ask_push_images(){

    # ask if images will be pushed to registry
    if [[ -z "${PUSH_IMAGES}" ]]; then

        echo -n "Will images be pushed to registry? (true/false): "
        read PUSH_IMAGES

    fi

    if ! [[ "$PUSH_IMAGES" == "true" || "$PUSH_IMAGES" == "false" ]]; then
        
        echo "Wrong input for \"PUSH_IMAGES\". (true/false)"
        exit 1;

    fi

}

function ask_contains_vdi(){

    # ask if image contains vdi
    if [[ -z "${CONTAINS_VDI}" ]]; then

        echo -n "Does image contain VDI? (true/false): "
        read CONTAINS_VDI

    fi

    if ! [[ "$CONTAINS_VDI" == "true" || "$CONTAINS_VDI" == "false" ]]; then
        
        echo "Wrong input for \"CONTAINS_VDI\". (true/false)"
        exit 1;

    fi

}

function ask_rebuild_robot(){
    
    # ask if robot will be rebuilt
    if [[ -z "${REBUILD_ROBOT}" ]]; then
        
        echo -n "Rebuild robot base? (true/false): "
        read REBUILD_ROBOT

    fi

    if ! [[ "$REBUILD_ROBOT" == "true" || "$REBUILD_ROBOT" == "false" ]]; then
        
        echo "Wrong input for \"REBUILD_ROBOT\". (true/false)"
        exit 1;

    fi

}

function ask_rebuild_vdi(){

    # ask if vdi will be rebuilt
    if [[ -z "${REBUILD_VDI}" ]]; then

        echo -n "Rebuild VDI? (true/false): "
        read REBUILD_VDI
    
    fi

    if ! [[ "$REBUILD_VDI" == "true" || "$REBUILD_VDI" == "false" ]]; then

        echo "Wrong input for \"REBUILD_VDI\". (true/false)"
        exit 1;

    fi
   
}

function docker_build_driver(){
    set -e

    docker build ./vdi/base -t $DRIVER_IMAGE \
    --build-arg NVIDIA_DRIVER_VERSION=$NVIDIA_DRIVER_VERSION \
    --build-arg GPU_AGNOSTIC=$GPU_AGNOSTIC \
    --build-arg BASE_IMAGE=$OPENGL_IMAGE

    docker_push_image $DRIVER_IMAGE
}

function docker_build_vdi(){
    set -e

    docker build ./vdi/$UBUNTU_DESKTOP -t $VDI_IMAGE \
    --build-arg BASE_IMAGE=$DRIVER_IMAGE

    docker tag $VDI_IMAGE $VDI_IMAGE-amd64
    docker_push_image $VDI_IMAGE
    docker_push_image $VDI_IMAGE-amd64
}

function docker_build_ros(){
    set -e

    docker buildx rm multiarch_builder || true
    docker buildx create --name multiarch_builder --use

    if [ "$MULTIPLE_ROS_DISTRO" == "true" ]; then

        SUB_ROS_IMAGE=sub-ros:$ROS_DISTRO_1
        
        docker buildx build ./ros/$ROS_DISTRO_1 -t $SUB_ROS_IMAGE \
            --build-arg BASE_IMAGE=$ROS_BASE_IMAGE \
            --platform linux/amd64,linux/arm64 \
            --push
        
        docker buildx build ./ros/$ROS_DISTRO_2 -t $ROS_IMAGE \
            --build-arg BASE_IMAGE=$SUB_ROS_IMAGE \
            --platform linux/amd64,linux/arm64 \
            --push
        
    else

        docker buildx build ./ros/$ROS_DISTRO -t $ROS_IMAGE \
            --build-arg BASE_IMAGE=$ROS_BASE_IMAGE \
            --platform linux/amd64,linux/arm64 \
            --push
       
    fi

}

function docker_build_robolaunch_robot(){
    set -e

    docker buildx rm multiarch_builder || true
    docker buildx create --name multiarch_builder --use

    if [ "$MULTIPLE_ROS_DISTRO" == "true" ]; then

        docker buildx build ./robot -t $ROBOLAUNCH_ROBOT_IMAGE \
        --build-arg BASE_IMAGE=$ROBOLAUNCH_ROBOT_BASE_IMAGE \
        --build-arg BRIDGE_DISTRO_1=$ROS_DISTRO_1 \
        --build-arg BRIDGE_DISTRO_2=$ROS_DISTRO_2 \
        --platform linux/amd64,linux/arm64 \
        --push

    else

        docker buildx build ./robot -t $ROBOLAUNCH_ROBOT_IMAGE \
        --build-arg BASE_IMAGE=$ROBOLAUNCH_ROBOT_BASE_IMAGE \
        --build-arg BRIDGE_DISTRO_1=$ROS_DISTRO \
        --build-arg BRIDGE_DISTRO_2=$ROS_DISTRO \
        --platform linux/amd64,linux/arm64 \
        --push

    fi
    
}

function docker_build_robot(){
    set -e

    docker buildx rm multiarch_builder || true
    docker buildx create --name multiarch_builder --use

    if [ "$MULTIPLE_ROS_DISTRO" == "true" ]; then

        docker buildx build ./robot/$ROBOT_NAME -t $ROBOT_IMAGE \
        --build-arg BASE_IMAGE=$ROBOT_BASE_IMAGE \
        --build-arg ROS_DISTRO=$ROBOT_DISTRO \
        --platform linux/amd64,linux/arm64 \
        --push

    else

        docker buildx build ./robot/$ROBOT_NAME -t $ROBOT_IMAGE \
        --build-arg BASE_IMAGE=$ROBOT_BASE_IMAGE \
        --build-arg ROS_DISTRO=$ROS_DISTRO \
        --platform linux/amd64,linux/arm64 \
        --push

    fi
    
}

function docker_push_image(){
    set -e
    
    if [ "$PUSH_IMAGES" == "true" ]; then
    
        docker push $1

    fi
}


cd "$(dirname "$0")"

ask_registry
ask_push_images
ask_contains_vdi
ask_rebuild_robot
ask_ubuntu_distro

if [ "$CONTAINS_VDI" == "true" ]; then
    
    if [ "$REBUILD_ROBOT" == "true" ]; then

        ask_rebuild_vdi

        if [ "$REBUILD_VDI" == "true" ]; then
            build_vdi
            build_ros
            build_robolaunch_robot
            build_robot
        else
            build_ros
            build_robolaunch_robot
            build_robot
        fi

    else

        build_robot

    fi

else

    if [ "$REBUILD_ROBOT" == "true" ]; then
        
        build_ros
        build_robolaunch_robot
        build_robot

    else

        build_robot

    fi

fi
