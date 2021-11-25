#!/usr/bin/env bash

CLUSTER_NAME="dvo-test"

if ! command -v k3d &> /dev/null
then
    echo "Installing k3d"
    curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
fi

if ! command -v kubectl &> /dev/null
then
    echo "Installing kubectl"
    curl -LO "https://dl.k8s.io/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl.sha256"
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
fi

if ! command -v helm &> /dev/null
then
    echo "installing helm"
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
fi

k3d cluster list | grep ${CLUSTER_NAME} &> /dev/null
SUCCESS=$?
if [ $SUCCESS -eq 1 ]
then
    echo "Creating ${CLUSTER_NAME} k3d cluster"
    k3d cluster create $CLUSTER_NAME -p "8080:30080@server:0" -p "8081:30081@server:0" --registry-create $CLUSTER_NAME-registry:0.0.0.0:5432 --registry-config ./registries.yaml
fi

echo ""

echo "Installing prometheus"
echo ""
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts &> /dev/null
helm install -n prom --create-namespace --set server.service.type=NodePort --set server.service.nodePort=30080 prom prometheus-community/prometheus &> /dev/null

echo "Installing fake prometheus metrics for test/verification"
kubectl apply -f fake-prom-exporter-deployment.yml
echo "Visit localhost:8080 and query \`fake_metric\`"
echo ""

echo "Installing deployment validation operator"
echo ""
kubectl create namespace deployment-validation-operator
kubectl apply -f ../../deploy/k8s/.
