#!/bin/bash

echo "Esperando a que Kibana esté listo..."

until curl -s http://kibana:5601/api/status | grep -q '"level":"available"'; do
  echo "Kibana no está listo aún, esperando..."
  sleep 5
done

echo "Kibana está listo. Configurando Data View..."

curl -X POST "http://kibana:5601/api/data_views/data_view" \
  -H "kbn-xsrf: true" \
  -H "Content-Type: application/json" \
  -d '{
    "data_view": {
      "title": "efact-logs-*",
      "name": "EFACT Logs",
      "timeFieldName": "@timestamp"
    }
  }'

echo ""
echo "Data View 'efact-logs-*' creado exitosamente"

curl -X POST "http://kibana:5601/api/data_views/data_view" \
  -H "kbn-xsrf: true" \
  -H "Content-Type: application/json" \
  -d '{
    "data_view": {
      "title": "ms1-logs-*",
      "name": "MS1 Logs",
      "timeFieldName": "@timestamp"
    }
  }'

echo ""
echo "Data View 'ms1-logs-*' creado exitosamente"

curl -X POST "http://kibana:5601/api/data_views/data_view" \
  -H "kbn-xsrf: true" \
  -H "Content-Type: application/json" \
  -d '{
    "data_view": {
      "title": "ms2-logs-*",
      "name": "MS2 Logs",
      "timeFieldName": "@timestamp"
    }
  }'

echo ""
echo "Data View 'ms2-logs-*' creado exitosamente"
