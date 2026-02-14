#!/bin/sh
set -e

if [ -z "$VAULT_ADDR" ] || [ -z "$VAULT_TOKEN" ]; then
  echo "Error: VAULT_ADDR y VAULT_TOKEN son requeridos"
  exit 1
fi

SPRING_PROFILES_ACTIVE=${SPRING_PROFILES_ACTIVE:-local}

echo "Esperando que cargue el Vault"
RETRY=0
until curl -s -f "$VAULT_ADDR/v1/sys/health" > /dev/null 2>&1; do
  RETRY=$((RETRY + 1))
  if [ $RETRY -ge 30 ]; then
    echo "Error: Vault no disponible"
    exit 1
  fi
  sleep 2
done

echo "Obteniendo secrets para ambente: $SPRING_PROFILES_ACTIVE"
SECRET_PATH="secret/data/efact/ms2/rsa-keys/$SPRING_PROFILES_ACTIVE"
RESPONSE=$(curl -s -H "X-Vault-Token: $VAULT_TOKEN" "$VAULT_ADDR/v1/$SECRET_PATH")

if ! echo "$RESPONSE" | jq -e '.data.data' > /dev/null 2>&1; then
  echo "Error: No se pudieron obtener secrets de Vault"
  exit 1
fi

export SIGNATURE_PRIVATE_KEY=$(echo "$RESPONSE" | jq -r '.data.data.private_key')
export SIGNATURE_PUBLIC_KEY=$(echo "$RESPONSE" | jq -r '.data.data.public_key')

if [ -z "$SIGNATURE_PRIVATE_KEY" ] || [ "$SIGNATURE_PRIVATE_KEY" = "null" ]; then
  echo "Error: No se pudo obtener private_key"
  exit 1
fi

if [ -z "$SIGNATURE_PUBLIC_KEY" ] || [ "$SIGNATURE_PUBLIC_KEY" = "null" ]; then
  echo "Error: No se pudo obtener public_key"
  exit 1
fi

echo "Claves cargadas"
exec java -jar app.jar
