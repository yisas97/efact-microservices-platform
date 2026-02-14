#!/bin/bash
set -e

ENVIRONMENT=${1:-dev}
NAMESPACE="efact"

[ "$ENVIRONMENT" == "dev" ] && NAMESPACE="efact-dev"

kubectl get namespace $NAMESPACE &> /dev/null || { echo "Namespace not found"; exit 1; }

forward_port() {
    kubectl port-forward -n $NAMESPACE svc/$1 $2:$3 &
}

forward_port "ms1-documents-service" 5000 5000
forward_port "ms2-validator-service" 8080 8080
forward_port "mongodb-service" 27017 27017
forward_port "rabbitmq-service" 5672 5672
forward_port "rabbitmq-service" 15672 15672
forward_port "elasticsearch-service" 9200 9200
forward_port "kibana-service" 5601 5601
forward_port "vault-service" 8200 8200

echo "Port-forward active. Press Ctrl+C to stop."
wait
