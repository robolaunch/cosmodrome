#!/bin/bash

# ORGANIZATION=rls-doo \
# TEAM=team-cloudy \
# REGION=eu-east-2 \
# CLOUD_INSTANCE=robot-cloud-02 \
# CLOUD_INSTANCE_ALIAS=instance-1 \
# PHYSICAL_INSTANCE=robot-cloudy-01 \
# DESIRED_CLUSTER_CIDR=10.20.1.0/24 \
# DESIRED_SERVICE_CIDR=10.20.2.0/24 \
# NETWORK=External \
# ./run.sh

set -e;

BLUE='\033[0;34m';
GREEN='\033[0;32m';
RED='\033[0;31m';
NC='\033[0m';

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

set_cloud_instance_ca () {
    if [[ -z "${CLOUD_INSTANCE_CA}" ]]; then
        print_err "Environment variable CLOUD_INSTANCE_CA should be set.";
    else
        CLOUD_INSTANCE_CA=$CLOUD_INSTANCE_CA;
    fi
}

set_cloud_instance_api_server () {
    if [[ -z "${CLOUD_INSTANCE_API_SERVER}" ]]; then
        print_err "Environment variable CLOUD_INSTANCE_API_SERVER should be set.";
    else
        CLOUD_INSTANCE_API_SERVER=$CLOUD_INSTANCE_API_SERVER;
    fi
}

set_cloud_instance_user () {
    if [[ -z "${CLOUD_INSTANCE_USER}" ]]; then
        print_err "Environment variable CLOUD_INSTANCE_USER should be set.";
    else
        CLOUD_INSTANCE_USER=$CLOUD_INSTANCE_USER;
    fi
}

set_cloud_instance_oauth_token () {
    if [[ -z "${CLOUD_INSTANCE_OAUTH_TOKEN}" ]]; then
        print_err "Environment variable CLOUD_INSTANCE_OAUTH_TOKEN should be set.";
    else
        CLOUD_INSTANCE_OAUTH_TOKEN=$CLOUD_INSTANCE_OAUTH_TOKEN;
    fi
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

set_physical_instance () {
    if [[ -z "${PHYSICAL_INSTANCE}" ]]; then
        print_err "Environment variable PHYSICAL_INSTANCE should be set.";
    else
        PHYSICAL_INSTANCE=$PHYSICAL_INSTANCE;
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

set_network () {
    if [[ -z "${NETWORK}" ]]; then
        print_err "Environment variable NETWORK should be set.";
    else
        NETWORK=$NETWORK;
    fi
}

set_connection_hub_key () {
    if [[ -z "${CONNECTION_HUB_KEY}" ]]; then
        print_err "Environment variable CONNECTION_HUB_KEY should be set.";
    else
        CONNECTION_HUB_KEY=$CONNECTION_HUB_KEY;
    fi
}

check_node_name () {
    NODE_NAME=$(kubectl get nodes -l node-role.kubernetes.io/master -o 'jsonpath={.items[*].metadata.name}');
}

check_cluster_cidr () {
    check_node_name
    PHYSICAL_INSTANCE_CLUSTER_CIDR=$(kubectl get nodes $NODE_NAME -o jsonpath='{.spec.podCIDR}');
}

check_service_cidr () {
    PHYSICAL_INSTANCE_SERVICE_CIDR=$(echo '{"apiVersion":"v1","kind":"Service","metadata":{"name":"tst"},"spec":{"clusterIP":"1.1.1.1","ports":[{"port":443}]}}' | kubectl apply -f - 2>&1 | sed 's/.*valid IPs is //');
}

check_inputs () {
    if [[ -z "${SKIP_PLATFORM}" ]]; then
        set_cloud_instance_ca;
        set_cloud_instance_api_server;
        set_cloud_instance_user;
        set_cloud_instance_oauth_token;
    fi
    set_organization;
    set_team;
    set_region;
    set_cloud_instance;
    set_cloud_instance_alias;
    set_physical_instance;
    set_connection_hub_key;
    set_desired_cluster_cidr;
    set_desired_service_cidr;
    set_network;
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

set_up_k3s () {
    print_log "Setting up k3s...";
    curl -sfL https://get.k3s.io | \
        INSTALL_K3S_VERSION=$K3S_VERSION+k3s1 \
        K3S_KUBECONFIG_MODE="644" \
        INSTALL_K3S_EXEC="  --cluster-cidr=$DESIRED_CLUSTER_CIDR --service-cidr=$DESIRED_SERVICE_CIDR    --cluster-domain=$PHYSICAL_INSTANCE.local --disable-network-policy --disable=traefik --disable=local-storage" sh -;
    sleep 5;
}

check_cluster () {
    check_cluster_cidr;
    check_service_cidr;
}

label_node () {
    print_log "Labeling node...";
    check_node_name;
    kubectl label --overwrite=true node $NODE_NAME \
        robolaunch.io/organization=$ORGANIZATION \
        robolaunch.io/region=$REGION \
        robolaunch.io/team=$TEAM \
        robolaunch.io/cloud-instance=$CLOUD_INSTANCE \
        robolaunch.io/cloud-instance-alias=$CLOUD_INSTANCE_ALIAS \
        robolaunch.io/physical-instance=$PHYSICAL_INSTANCE \
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
    helm upgrade -i openebs openebs/openebs \
    --namespace openebs \
    --create-namespace;
    sleep 5;
    kubectl patch storageclass openebs-hostpath -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}';
}

install_cert_manager () {
    print_log "Installing cert-manager... This might take around one minute.";
    kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/$CERT_MANAGER_VERSION/cert-manager.yaml;
    # TODO: Check if cert-manager is up & running.
    sleep 30;
}

install_operator_suite () {
    print_log "Installing operator Helm charts... This might take around one minute."
    
    CHO_HELM_INSTALL_SUCCEEDED="false"
    while [ "$CHO_HELM_INSTALL_SUCCEEDED" != "true" ]
    do 
        CHO_HELM_INSTALL_SUCCEEDED="true"
        helm upgrade -i \
            connection-hub-operator robolaunch/connection-hub-operator \
            --namespace connection-hub-system \
            --create-namespace \
            --version $CONNECTION_HUB_OPERATOR_CHART_VERSION || CHO_HELM_INSTALL_SUCCEEDED="false";
        sleep 1;
    done

    RO_HELM_INSTALL_SUCCEEDED="false"
    while [ "$RO_HELM_INSTALL_SUCCEEDED" != "true" ]
    do 
        RO_HELM_INSTALL_SUCCEEDED="true"
        helm upgrade -i \
            robot-operator robolaunch/robot-operator \
            --namespace robot-system \
            --create-namespace \
            --version $ROBOT_OPERATOR_CHART_VERSION || RO_HELM_INSTALL_SUCCEEDED="false";
        sleep 1;
    done

    FO_HELM_INSTALL_SUCCEEDED="false"
    while [ "$FO_HELM_INSTALL_SUCCEEDED" != "true" ]
    do 
        FO_HELM_INSTALL_SUCCEEDED="true"
        helm upgrade -i \
            fleet-operator robolaunch/fleet-operator \
            --namespace fleet-system \
            --create-namespace \
            --version $FLEET_OPERATOR_CHART_VERSION || FO_HELM_INSTALL_SUCCEEDED="false";
        sleep 1;
    done

    sleep 15;
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

join_connection_hub () {
    print_log "Joining connection hub...";
    check_cluster
    echo $CONNECTION_HUB_KEY | base64 --decode > ch-pi.yaml;
    yq e -i ".metadata.labels.\"robolaunch.io/physical-instance\" = \"$PHYSICAL_INSTANCE\"" ch-pi.yaml;
    yq e -i ".spec.submarinerSpec.clusterCIDR = \"$PHYSICAL_INSTANCE_CLUSTER_CIDR\"" ch-pi.yaml;
    yq e -i ".spec.submarinerSpec.serviceCIDR = \"$PHYSICAL_INSTANCE_SERVICE_CIDR\"" ch-pi.yaml;
    yq e -i ".spec.submarinerSpec.networkType = \"$NETWORK\"" ch-pi.yaml;
    
    CH_INSTALL_SUCCEEDED="false"
    while [ "$CH_INSTALL_SUCCEEDED" != "true" ]
    do 
        CH_INSTALL_SUCCEEDED="true"
        kubectl apply -f ch-pi.yaml || CH_INSTALL_SUCCEEDED="false";
        sleep 3;
    done
    
    rm -rf ch-pi.yaml;
    check_connection_hub_phase;
}

check_cloud_instance_phase () {
    while [ true ]
    do
        CLOUD_INSTANCE_PHASE=$(kubectl get cloudinstances $CLOUD_INSTANCE -o jsonpath="{.status.phase}" | yq -P);
        if [ "$CLOUD_INSTANCE_PHASE" = "Connected" ]; then
            break;
        fi

        print_log "Checking connection status...";
        sleep 3;
    done
}

register_me () {
    check_cluster;
    PHYSICAL_INSTANCE_API_SERVER_URL="https://"${PHYSICAL_INSTANCE_CLUSTER_CIDR%0/*}"1:6443"
    CERT_AUTHORITY_DATA=$(yq '.clusters[] | select(.name == "default") | .cluster.certificate-authority-data' $KUBECONFIG);
    CLIENT_CERTIFICATE=$(yq '.users[] | select(.name == "default") | .user.client-certificate-data' $KUBECONFIG);
    CLIENT_KEY=$(yq '.users[] | select(.name == "default") | .user.client-key-data' $KUBECONFIG);
    
    # SKIP_PLATFORM is not set
    if [[ -z "${SKIP_PLATFORM}" ]]; then
        echo $CLOUD_INSTANCE_CA | base64 --decode >> /tmp/ca.crt;
        kubectl config set-cluster cloud-instance --server=$CLOUD_INSTANCE_API_SERVER --certificate-authority=/tmp/ca.crt --embed-certs=true 2>/dev/null 1>/dev/null;
        kubectl config set-credentials $CLOUD_INSTANCE_USER --token=$CLOUD_INSTANCE_OAUTH_TOKEN 2>/dev/null 1>/dev/null;
        kubectl config set-context $CLOUD_INSTANCE_USER --cluster=cloud-instance --user=$CLOUD_INSTANCE_USER 2>/dev/null 1>/dev/null;
        kubectl config use-context $CLOUD_INSTANCE_USER 2>/dev/null 1>/dev/null;
        sleep 2;
        cat << EOF | kubectl apply -f -
apiVersion: connection-hub.roboscale.io/v1alpha1
kind: PhysicalInstance
metadata:
  name: $PHYSICAL_INSTANCE
spec:
  server: $PHYSICAL_INSTANCE_API_SERVER_URL
  credentials:
    certificateAuthority: $CERT_AUTHORITY_DATA
    clientCertificate: $CLIENT_CERTIFICATE
    clientKey: $CLIENT_KEY
EOF
        print_log "Go to the platform and check if your physical instance $PHYSICAL_INSTANCE is connected to your cloud instance $CLOUD_INSTANCE_ALIAS/$CLOUD_INSTANCE.";
        print_log "Check your physical instance status by running the command below in your cloud instance:\n\n";
        printf "\n\n";
        printf "watch kubectl get physicalinstances";
        printf "\n\n";
    # SKIP_PLATFORM is set
    else
        printf "\n\n"
        echo \
"cat <<EOF | kubectl apply -f -
apiVersion: connection-hub.roboscale.io/v1alpha1
kind: PhysicalInstance
metadata:
  name: $PHYSICAL_INSTANCE
spec:
  server: $PHYSICAL_INSTANCE_API_SERVER_URL
  credentials:
    certificateAuthority: $CERT_AUTHORITY_DATA
    clientCertificate: $CLIENT_CERTIFICATE
    clientKey: $CLIENT_KEY
EOF";
        printf "\n";
        print_log "Physical instance is connected to the cloud instance $CLOUD_INSTANCE_ALIAS/$CLOUD_INSTANCE.";
        print_log "In order to complete registration of physical instance you should run the command above in cloud instance $CLOUD_INSTANCE_ALIAS/$CLOUD_INSTANCE.";
    fi
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
K3S_VERSION=v$(yq ''"${VERSION_SELECTOR_STR}"' | .roboticsCloud.kubernetes.version' < platform.yaml)
CERT_MANAGER_VERSION=$(yq ''"${VERSION_SELECTOR_STR}"' | .roboticsCloud.kubernetes.components.cert-manager.version' < platform.yaml)
CONNECTION_HUB_OPERATOR_CHART_VERSION=$(yq ''"${VERSION_SELECTOR_STR}"' | .roboticsCloud.kubernetes.operators.connectionHub.helm.version' < platform.yaml)
ROBOT_OPERATOR_CHART_VERSION=$(yq ''"${VERSION_SELECTOR_STR}"' | .roboticsCloud.kubernetes.operators.robot.helm.version' < platform.yaml)
FLEET_OPERATOR_CHART_VERSION=$(yq ''"${VERSION_SELECTOR_STR}"' | .roboticsCloud.kubernetes.operators.fleet.helm.version' < platform.yaml)

opening >&3
(check_inputs)

print_global_log "Installing tools...";
(install_post_tools)

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

print_global_log "Installing cert-manager...";
(install_cert_manager)

print_global_log "Installing robolaunch Operator Suite...";
(install_operator_suite)

print_global_log "Joining connection hub $CLOUD_INSTANCE_ALIAS/$CLOUD_INSTANCE...";
(join_connection_hub)

print_global_log "Checking connection status...";
(check_cloud_instance_phase)

register_me >&3