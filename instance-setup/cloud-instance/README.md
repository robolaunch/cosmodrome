# Setting Up Cloud Instance (VM)

This is the table of contents for this document.

- [Quick Start](#quick-start)
  - [Getting Connection Hub Key](#get-connection-hub-key)
- [Components](#components)

## Quick Start

This document assumes you have a virtual machine (with `amd64` architecture) provisioned on any cloud provider and has Ubuntu 20.04 or 22.04 on it.

Run the following command to set up cloud instance:

```bash
ORGANIZATION=sample-org \
TEAM=sample-team \
REGION=sample-region \
CLOUD_INSTANCE=cloud-instance \
CLOUD_INSTANCE_ALIAS=my-first-instance \
DESIRED_CLUSTER_CIDR=10.100.1.0/24 \
DESIRED_SERVICE_CIDR=10.100.2.0/24 \
curl https://raw.githubusercontent.com/robolaunch/cosmodrome/main/instance-setup/cloud-instance/run.sh | bash
```

The output will be similar to this:
```
Installing tools...
Setting up NVIDIA container runtime...
Setting up k3s cluster...
Checking cluster health...
Labeling node...
Updating Helm repositories...
Installing openebs...
Installing NVIDIA device plugin...
Installing cert-manager...
Installing robolaunch Operator Suite...
Deploying Connection Hub...


Get connection hub key by running the command below:

kubectl get connectionhub connection-hub -o jsonpath={.status.key} | yq -P

You can use this key to establish a connection with cloud instance my-first-instance/cloud-instance.
```

It's done. Now you can get the connection hub key in the next section.

### Getting Connection Hub Key

```bash
kubectl get connectionhub connection-hub -o jsonpath={.status.key} | yq -P
```

You should save this key since it will be used to register physical instances.

## Components