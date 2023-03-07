#!/bin/bash

# prerequisites: kubectl, helm, curl, wget, yq, figlet

set -e;

figlet 'robolaunch' -f slant;
printf "\n";
echo "\"We Empower ROS/ROS2 based GPU Offloaded Robots & Geographically Distributed Fleets\"";
printf "\n";

BLUE='\033[0;34m';
RED='\033[0;31m';
NC='\033[0m';

CERT_MANAGER_VERSION="v1.8.0";
OPERATOR_SUITE_VERSION="0.1.0";
CH_CLOUD_INSTANCE_URL="https://gist.githubusercontent.com/tunahanertekin/f041e2c3fbc6cdaadd72816c350b357c/raw/ac86a73e70ea8dce5903eed3472b26afdc255f0d/ch-ci.yaml";

print_log () {
    echo -e "${BLUE}$1${NC}";
}

print_err () {
    echo -e "${RED}Error: $1${NC}";
    exit 1;
}

preflight_checks () {
    set_organization;
    set_team;
    set_region;
    set_cloud_instance;
    set_cloud_instance_alias;
    set_api_server_url;
    set_cluster_cidr;
    set_service_cidr;
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
        set_cluster_root_domain;
        REAL_CLOUD_INSTANCE=${CLUSTER_ROOT_DOMAIN%.local*};
        
        if [ "$CLOUD_INSTANCE" != "$REAL_CLOUD_INSTANCE" ]; then
            print_err "Cloud instance name should be derived from cluster root domain.";
        fi
    fi
}

set_cloud_instance_alias () {
    if [[ -z "${CLOUD_INSTANCE_ALIAS}" ]]; then
        print_err "Environment variable CLOUD_INSTANCE_ALIAS should be set.";
    else
        CLOUD_INSTANCE_ALIAS=$CLOUD_INSTANCE_ALIAS;
    fi
}

set_public_ip () {
    if [[ -z "${PUBLIC_IP}" ]]; then
        PUBLIC_IP=$(curl https://ipinfo.io/ip);
    else
        PUBLIC_IP=$PUBLIC_IP;
    fi
}

set_api_server_url () {
    set_public_ip
    CLOUD_INSTANCE_API_SERVER_URL="$PUBLIC_IP:6443";
}

set_node_name () {
    NODE_NAME=$(kubectl get nodes -l node-role.kubernetes.io/master -o 'jsonpath={.items[*].metadata.name}');
}

set_cluster_cidr () {
    set_node_name;
    CLOUD_INSTANCE_CLUSTER_CIDR=$(kubectl get nodes $NODE_NAME -o jsonpath='{.spec.podCIDR}');
}

set_service_cidr () {
    CLOUD_INSTANCE_SERVICE_CIDR=$(echo '{"apiVersion":"v1","kind":"Service","metadata":{"name":"tst"},"spec":{"clusterIP":"1.1.1.1","ports":[{"port":443}]}}' | kubectl apply -f - 2>&1 | sed 's/.*valid IPs is //');
}

set_up_k3s () {
    print_log "Setting up k3s...";
}

install_helm () {
    print_log "Installing Helm...";
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash;
}

update_helm_repositories () {
    print_log "Updating Helm repositories...";
    helm repo add jetstack https://charts.jetstack.io;
    helm repo add robolaunch http://charts.robolaunch.dev/helm;
    helm repo update;
}

install_cert_manager () {
    print_log "Installing cert-manager... This might take around one minute.";
    kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/$CERT_MANAGER_VERSION/cert-manager.yaml;
    # TODO: Check if cert-manager is up & running.
    sleep 30;
}

install_helm_charts () {
    print_log "Installing operator Helm charts... This might take around one minute."
    HELM_INSTALL_SUCCEEDED="false"

    while [ "$HELM_INSTALL_SUCCEEDED" != "true" ]
    do 
        HELM_INSTALL_SUCCEEDED="true"
        helm upgrade -i \
            operator-suite robolaunch/operator-suite \
            --set global.organization=$ORGANIZATION \
            --set global.team=$TEAM \
            --set global.region=$REGION \
            --set global.cloudInstance=$CLOUD_INSTANCE \
            --set global.cloudInstanceAlias=$CLOUD_INSTANCE_ALIAS \
            --version $OPERATOR_SUITE_VERSION || HELM_INSTALL_SUCCEEDED="false";
    done

    sleep 30;
}

check_connection_hub_phase () {
    while [ true ]
    do
        CH_PHASE=$(kubectl get connectionhub connection-hub -o jsonpath=\"{.status.phase}\" | yq -P);
        if [ "$CH_PHASE" = "ReadyForOperation" ]; then
            print_log "Connection hub is ready to establish connections.";
            break;
        fi

        print_log "Checking connection hub phase -> $CH_PHASE";
        sleep 3;
    done
}

deploy_connection_hub () {
    print_log "Deploying connection hub..."
    wget $CH_CLOUD_INSTANCE_URL;
    yq e -i ".metadata.labels.\"robolaunch.io/cloud-instance\" = \"$CLOUD_INSTANCE\"" ch-ci.yaml;
    yq e -i ".metadata.labels.\"robolaunch.io/cloud-instance-alias\" = \"$CLOUD_INSTANCE_ALIAS\"" ch-ci.yaml;
    yq e -i ".spec.submarinerSpec.apiServerURL = \"$CLOUD_INSTANCE_API_SERVER_URL\"" ch-ci.yaml;
    yq e -i ".spec.submarinerSpec.clusterCIDR = \"$CLOUD_INSTANCE_CLUSTER_CIDR\"" ch-ci.yaml;
    yq e -i ".spec.submarinerSpec.serviceCIDR = \"$CLOUD_INSTANCE_SERVICE_CIDR\"" ch-ci.yaml;
    
    CH_INSTALL_SUCCEEDED="false"
    while [ "$CH_INSTALL_SUCCEEDED" != "true" ]
    do 
        CH_INSTALL_SUCCEEDED="true"
        kubectl apply -f ch-ci.yaml || CH_INSTALL_SUCCEEDED="false";
    done

    check_connection_hub_phase;
}

display_connection_hub_key () {
    CONNECTION_HUB_KEY=$(kubectl get connectionhub connection-hub -o jsonpath="{.status.key}" | yq -P);
    print_log "You can use this key to establish a connection with cloud instance $CLOUD_INSTANCE_ALIAS/$CLOUD_INSTANCE";
    printf "\n";
    echo $CONNECTION_HUB_KEY;
}

preflight_checks
# set_up_k3s
# install_helm
install_cert_manager
update_helm_repositories
install_helm_charts
deploy_connection_hub
display_connection_hub_key