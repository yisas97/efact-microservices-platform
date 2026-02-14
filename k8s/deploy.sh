#!/bin/bash
set -e

ENVIRONMENT=${1:-dev}
NAMESPACE="efact"

[ "$ENVIRONMENT" == "dev" ] && NAMESPACE="efact-dev"

command -v kubectl &> /dev/null || { echo "kubectl not installed"; exit 1; }
kubectl cluster-info &> /dev/null || { echo "Cannot connect to cluster"; exit 1; }

kubectl get namespace $NAMESPACE &> /dev/null || kubectl create namespace $NAMESPACE

kubectl apply -k overlays/$ENVIRONMENT/

kubectl wait --for=condition=ready pod -l app=mongodb -n $NAMESPACE --timeout=300s || true
kubectl wait --for=condition=ready pod -l app=vault -n $NAMESPACE --timeout=300s || true

echo ""
kubectl get pods -n $NAMESPACE
echo ""
kubectl get svc -n $NAMESPACE
