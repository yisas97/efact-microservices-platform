#!/bin/bash
set -e

ENVIRONMENT=${1:-dev}
NAMESPACE="efact"

[ "$ENVIRONMENT" == "dev" ] && NAMESPACE="efact-dev"

kubectl get namespace $NAMESPACE &> /dev/null || { echo "Namespace $NAMESPACE not found"; exit 1; }

echo "PODS"
kubectl get pods -n $NAMESPACE -o wide

echo ""
echo "SERVICES"
kubectl get svc -n $NAMESPACE

echo ""
echo "PVCs"
kubectl get pvc -n $NAMESPACE

echo ""
echo "DEPLOYMENTS"
kubectl get deployments -n $NAMESPACE

echo ""
echo "STATEFULSETS"
kubectl get statefulsets -n $NAMESPACE
