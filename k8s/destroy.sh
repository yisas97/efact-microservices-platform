#!/bin/bash
set -e

ENVIRONMENT=${1:-dev}
KEEP_DATA=${2:-}
NAMESPACE="efact"

[ "$ENVIRONMENT" == "dev" ] && NAMESPACE="efact-dev"

echo "Deleting $NAMESPACE namespace"
read -p "Confirm (yes/no): " -r
[[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]] && exit 1

kubectl delete -k overlays/$ENVIRONMENT/ --ignore-not-found=true

[ "$KEEP_DATA" != "--keep-data" ] && kubectl delete pvc --all -n $NAMESPACE --ignore-not-found=true

kubectl delete namespace $NAMESPACE --ignore-not-found=true
