# SAP Credential Store Provider for Secrets Store CSI Driver 

SAP Credential Store provider for the [Secrets Store CSI driver](https://secrets-store-csi-driver.sigs.k8s.io/) enables you to pull passwords and encryption keys from the [SAP Credential Store](https://help.sap.com/docs/credential-store/sap-credential-store/what-is-sap-credential-store) and mount them into Kubernetes pods.

## Installation

You can use the deployment manifests in [deploy](./deploy/) which install the following components:

* Secrets Store CSI driver via the [Helm chart](https://artifacthub.io/packages/helm/secret-store-csi-driver/secrets-store-csi-driver)
* SAP Credential Store provider

You can install them in your current Kubernetes cluster by using [Kustomize](https://kustomize.io/):

```shell
kubectl kustomize --enable-helm deploy/ | kubectl apply -f-
```

Note: This provider requires an mTLS service key to communicate with the SAP Credentials Store (placed in [service-key.json](./deploy/service-key.json)). Check [this documentation link](https://help.sap.com/docs/credential-store/sap-credential-store/create-download-and-delete-service-key) which explains how to create one. [The SAP BTP Service Operator](https://github.com/SAP/sap-btp-service-operator) can also be used for automatic creation and rotation of such service keys.

### Usage

These [example manifests](./example/) demonstrate the basic scenario of mounting a password and key credentials into a pod.

The credential's metadata is described in the [secret-provider-class.yaml](./example/secret-provider-class.yaml) and follows this syntax:

* `name` - name of the source credential in SAP Credential Store
* `namespace` - namespace of the source credential in SAP Credential Store
* `type` - type of the source credential in SAP Credential Store, either *key* or *password*
* `fileName` - name of the destination file which will be mounted in the K8s pod
* `mode` - permissions of the destination file, e.g., *0640*, *0400*, *0777*. Defaults to *0644* if omitted

### Local Setup

```shell
# Build a custom container image
make image

# Setup local K8s cluster with kind: https://kind.sigs.k8s.io/
# The command also deploys the provider and the secrets store csi driver
make setup-kind
kubectl get pod -n csi
```