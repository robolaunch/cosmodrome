#!/bin/bash

# ORGANIZATION=rls-doo \
# TEAM=team-nlp \
# REGION=eu-east-2 \
# CLOUD_INSTANCE=dev-cloud-02 \
# CLOUD_INSTANCE_ALIAS=instance-1 \
# DESIRED_CLUSTER_CIDR=10.100.1.0/24 \
# DESIRED_SERVICE_CIDR=10.100.2.0/24 \
# ./run.sh

set -e;

BLUE='\033[0;34m';
GREEN='\033[0;32m';
RED='\033[0;31m';
NC='\033[0m';

CH_CLOUD_INSTANCE_URL="https://gist.githubusercontent.com/tunahanertekin/f041e2c3fbc6cdaadd72816c350b357c/raw/ac86a73e70ea8dce5903eed3472b26afdc255f0d/ch-ci.yaml";
ARCH=$(dpkg --print-architecture)
TIMESTAMP=$(date +%s)
OUTPUT_FILE="out_$TIMESTAMP.log"

export KUBECONFIG="/etc/rancher/k3s/k3s.yaml";
exec 3>&1 >$OUTPUT_FILE 2>&1;

print_global_log () {
    echo -e "${GREEN}$1${NC}" >&3;
}

print_log () {
    echo -e "${GREEN}$1${NC}";
}

print_err () {
    echo -e "${RED}Error: $1${NC}" >&3;
    exit 1;
}

set_cluster_root_domain () {
    CLUSTER_ROOT_DOMAIN=$(kubectl get cm coredns -n kube-system -o jsonpath="{.data.Corefile}" \
        | grep ".local " \
        | awk -F ' ' '{print $2}');
}

set_organization () {
    if [[ -z "${ORGANIZATION}" ]]; then
        print_err "Environment variable ORGANIZATION should be set.";
    else
        ORGANIZATION=$ORGANIZATION;
    fi
}

set_team () {
    if [[ -z "${TEAM}" ]]; then
        print_err "Environment variable TEAM should be set.";
    else
        TEAM=$TEAM;
    fi
}

set_region () {
    if [[ -z "${REGION}" ]]; then
        print_err "Environment variable REGION should be set.";
    else
        REGION=$REGION;
    fi
}

set_cloud_instance () {
    if [[ -z "${CLOUD_INSTANCE}" ]]; then
        print_err "Environment variable CLOUD_INSTANCE should be set.";
    else
        CLOUD_INSTANCE=$CLOUD_INSTANCE;
    fi
}

set_cloud_instance_alias () {
    if [[ -z "${CLOUD_INSTANCE_ALIAS}" ]]; then
        print_err "Environment variable CLOUD_INSTANCE_ALIAS should be set.";
    else
        CLOUD_INSTANCE_ALIAS=$CLOUD_INSTANCE_ALIAS;
    fi
}

set_desired_cluster_cidr () {
    if [[ -z "${DESIRED_CLUSTER_CIDR}" ]]; then
        print_err "Environment variable DESIRED_CLUSTER_CIDR should be set.";
    else
        DESIRED_CLUSTER_CIDR=$DESIRED_CLUSTER_CIDR;
    fi
}

set_desired_service_cidr () {
    if [[ -z "${DESIRED_SERVICE_CIDR}" ]]; then
        print_err "Environment variable DESIRED_SERVICE_CIDR should be set.";
    else
        DESIRED_SERVICE_CIDR=$DESIRED_SERVICE_CIDR;
    fi
}

set_public_ip () {
    if [[ -z "${PUBLIC_IP}" ]]; then
        PUBLIC_IP=$(curl https://ipinfo.io/ip);
    else
        PUBLIC_IP=$PUBLIC_IP;
    fi
}

check_api_server_url () {
    set_public_ip
    CLOUD_INSTANCE_API_SERVER_URL="$PUBLIC_IP:6443";
}

check_node_name () {
    NODE_NAME=$(kubectl get nodes -l node-role.kubernetes.io/master -o 'jsonpath={.items[*].metadata.name}');
}

check_cluster_cidr () {
    check_node_name;
    CLOUD_INSTANCE_CLUSTER_CIDR=$(kubectl get nodes $NODE_NAME -o jsonpath='{.spec.podCIDR}');
}

check_service_cidr () {
    CLOUD_INSTANCE_SERVICE_CIDR=$(echo '{"apiVersion":"v1","kind":"Service","metadata":{"name":"tst"},"spec":{"clusterIP":"1.1.1.1","ports":[{"port":443}]}}' | kubectl apply -f - 2>&1 | sed 's/.*valid IPs is //');
}

check_inputs () {
    set_organization;
    set_team;
    set_region;
    set_cloud_instance;
    set_cloud_instance_alias;
    set_desired_cluster_cidr;
    set_desired_service_cidr;
}

get_versioning_map () {
    wget https://raw.githubusercontent.com/robolaunch/robolaunch/main/platform.yaml;
}

opening () {
    apt-get update 2>/dev/null 1>/dev/null;
    apt-get install -y figlet 2>/dev/null 1>/dev/null; 
    figlet 'robolaunch' -f slant;
    echo "Cloud Robotics Platform $PLATFORM_VERSION";
    printf "\n";
    echo "\"We Empower ROS/ROS2 based GPU Offloaded Robots & Geographically Distributed Fleets\"";
    printf "\n";
}

check_if_root () {
    if [ $USER != "root" ]; then
        print_err "You should switch to root using \"sudo -i\" before setup."
    fi
}

install_pre_tools () {
    print_log "Installing Tools...";
    # apt packages
    apt-get update;
    apt-get install -y curl wget;
    # install yq
    wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_${ARCH};
    chmod a+x /usr/local/bin/yq;
}

install_post_tools () {
    print_log "Installing Tools...";
    # helm
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash;
    # install kubectl
    curl -LO https://dl.k8s.io/release/$K3S_VERSION/bin/linux/${ARCH}/kubectl;
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl;
    rm -rf kubectl;
}

set_up_nvidia_container_runtime () {
    print_log "Setting up NVIDIA container runtime...";
    DEBIAN_FRONTEND=noninteractive
    apt-get update;
    apt-get install -y gnupg linux-headers-$(uname -r);
    apt-get install -y --no-install-recommends nvidia-driver-470;
    distribution=$(. /etc/os-release;echo $ID$VERSION_ID);
    curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | apt-key add -;
    curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | tee /etc/apt/sources.list.d/nvidia-docker.list;
    apt-get update;
    apt-get install -y nvidia-container-runtime;
}

set_up_k3s () {
    print_log "Setting up k3s...";
    curl -sfL https://get.k3s.io | \
        INSTALL_K3S_VERSION=$K3S_VERSION+k3s1 \
        K3S_KUBECONFIG_MODE="644" \
        INSTALL_K3S_EXEC="  --cluster-cidr=$DESIRED_CLUSTER_CIDR --service-cidr=$DESIRED_SERVICE_CIDR    --cluster-domain=$CLOUD_INSTANCE.local --disable-network-policy --disable=traefik --disable=local-storage" sh -;
    sleep 5;
}

check_cluster () {
    print_log "Checking cluster health...";
    check_api_server_url;
    check_cluster_cidr;
    check_service_cidr;
    set_public_ip
    curl -vk --resolve $PUBLIC_IP:6443:127.0.0.1  https://$PUBLIC_IP:6443/ping;
}

label_node () {
    print_log "Labeling node...";
    check_node_name;
    kubectl label --overwrite=true node $NODE_NAME \
        robolaunch.io/platform=$PLATFORM_VERSION \
        robolaunch.io/organization=$ORGANIZATION \
        robolaunch.io/region=$REGION \
        robolaunch.io/team=$TEAM \
        robolaunch.io/cloud-instance=$CLOUD_INSTANCE \
        robolaunch.io/cloud-instance-alias=$CLOUD_INSTANCE_ALIAS \
        submariner.io/gateway="true";
}

update_helm_repositories () {
    print_log "Updating Helm repositories...";
    helm repo add openebs https://openebs.github.io/charts;
    helm repo add jetstack https://charts.jetstack.io;
    helm repo add robolaunch https://robolaunch.github.io/charts;
    helm repo update;
}

install_openebs () {
    print_log "Installing openebs... This might take around one minute.";
    helm install openebs openebs/openebs \
    --namespace openebs \
    --create-namespace;
    sleep 5;
    kubectl patch storageclass openebs-hostpath -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}';
}

install_nvidia_runtime_class () {
    print_log "Installing NVIDIA runtime class...";
    cat << EOF | kubectl apply -f -
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: nvidia
handler: nvidia
EOF
}

install_nvidia_device_plugin () {
    print_log "Installing nvidia-device-plugin... This might take around one minute.";
    echo "version: v1
sharing:
  timeSlicing:
    resources:
    - name: nvidia.com/gpu
      replicas: 4 # number of slice for 1 core" > nvidia-device-plugin-config.yaml
    wget https://github.com/robolaunch/k8s-device-plugin/releases/download/v0.13.0/nvidia-device-plugin-0.13.0.tgz;
    helm upgrade -i nvdp ./nvidia-device-plugin-0.13.0.tgz \
    --version=0.13.0 \
    --namespace nvidia-device-plugin \
    --create-namespace \
    --set-file config.map.config=nvidia-device-plugin-config.yaml \
    --set runtimeClassName=nvidia;
    rm -rf nvidia-device-plugin-0.13.0.tgz;
    rm -rf nvidia-device-plugin-config.yaml;
}

install_cert_manager () {
    print_log "Installing cert-manager... This might take around one minute.";
    kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/$CERT_MANAGER_VERSION/cert-manager.yaml;
    # TODO: Check if cert-manager is up & running.
    sleep 30;
}

install_operator_suite () {
    print_log "Installing DevSpace Operator Helm chart... This might take around one minute."

    DO_HELM_INSTALL_SUCCEEDED="false"
    while [ "$DO_HELM_INSTALL_SUCCEEDED" != "true" ]
    do 
        DO_HELM_INSTALL_SUCCEEDED="true"
        helm upgrade -i \
            devspace-operator robolaunch/devspace-operator \
            --namespace devspace-system \
            --create-namespace \
            --version $DEVSPACE_OPERATOR_CHART_VERSION || DO_HELM_INSTALL_SUCCEEDED="false";
        sleep 1;
    done

    sleep 15;
}

display_ending_msg () {
    print_log "robolaunch DevCloud setup is finished successfully! You can now operate your development environments inside $CLOUD_INSTANCE_ALIAS/$CLOUD_INSTANCE.";
}

print_global_log "Waiting for the preflight checks...";
(check_if_root)
(install_pre_tools)
(get_versioning_map)

# Specifying platform & component versions
if [[ -z "${PLATFORM_VERSION}" ]]; then
    PLATFORM_VERSION=$(yq '.versions[0].version' < platform.yaml)
fi
VERSION_SELECTOR_STR='.versions[] | select(.version == "'"$PLATFORM_VERSION"'")'
K3S_VERSION=v$(yq ''"${VERSION_SELECTOR_STR}"' | .devCloud.kubernetes.version' < platform.yaml)
CERT_MANAGER_VERSION=$(yq ''"${VERSION_SELECTOR_STR}"' | .devCloud.kubernetes.components.cert-manager.version' < platform.yaml)
DEVSPACE_OPERATOR_CHART_VERSION=$(yq ''"${VERSION_SELECTOR_STR}"' | .devCloud.kubernetes.operators.devspace.helm.version' < platform.yaml)

opening >&3
(check_inputs)

print_global_log "Installing tools...";
(install_post_tools)

print_global_log "Setting up NVIDIA container runtime...";
(set_up_nvidia_container_runtime)

print_global_log "Setting up k3s cluster...";
(set_up_k3s)

print_global_log "Checking cluster health...";
(check_cluster)

print_global_log "Labeling node...";
(label_node)

print_global_log "Updating Helm repositories...";
(update_helm_repositories)

print_global_log "Installing openebs...";
(install_openebs)

print_global_log "Installing NVIDIA runtime class...";
(install_nvidia_runtime_class)

print_global_log "Installing NVIDIA device plugin...";
(install_nvidia_device_plugin)

print_global_log "Installing cert-manager...";
(install_cert_manager)

print_global_log "Installing robolaunch DevSpace Operator...";
(install_devspace_operator)

display_ending_msg >&3