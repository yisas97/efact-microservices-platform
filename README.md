# Prueba Técnica EFACT - Backend

Sistema de microservicios para gestión, validación y firma digital de documentos electrónicos.

## Índice

- [Definición](#definición)
- [Componentes](#componentes)
  - [ms1-documents](#ms1-documents)
  - [ms2-validator](#ms2-validator)
  - [MongoDB](#mongodb)
  - [RabbitMQ](#rabbitmq)
  - [Vault](#vault)
- [Arquitectura](#arquitectura)
- [Requisitos](#requisitos)
- [Levantar el Sistema](#levantar-el-sistema)
- [Portales de Acceso](#portales-de-acceso)
  - [API y Documentación](#api-y-documentación)
  - [Gestión de Infraestructura](#gestión-de-infraestructura)
  - [Monitoreo y Logs](#monitoreo-y-logs)
- [Validaciones](#validaciones)
- [Ejemplo de Uso](#ejemplo-de-uso)
- [Detener](#detener)

## Definición

Plataforma distribuida que permite crear documentos fiscales, validar cálculos de IGV y generar firmas digitales RSA mediante procesamiento asíncrono. Los documentos son validados automáticamente y firmados si cumplen con las reglas de negocio.

## Componentes

### ms1-documents
API REST desarrollada en Go con framework Gin. Gestiona el ciclo de vida completo de documentos (CRUD) y publica eventos a RabbitMQ.

Endpoints principales:
- POST /documents - Crear documento
- GET /documents - Listar documentos
- GET /documents/:id - Obtener por ID
- PUT /documents/:id - Actualizar documento
- DELETE /documents/:id - Eliminar documento

### ms2-validator
Servicio de validación en Java Spring Boot. Consume mensajes de RabbitMQ, valida cálculos de IGV (18%), genera firmas digitales RSA 2048 bits y actualiza documentos.

### MongoDB
Base de datos de documentos. Puerto 27017. Almacena documentos con índice único en idDocumento.

### RabbitMQ
Sistema de mensajería asíncrona. Puerto 5672. Cola principal: documents.created.

### Vault
Gestor de secrets. Puerto 8200. Almacena claves RSA privadas para firma digital.

## Arquitectura

```
Cliente → MS1 (API REST) → MongoDB
              ↓
          RabbitMQ
              ↓
         MS2 (Validator) → MongoDB
              ↓
            Vault
```

Flujo:
1. Cliente crea documento en MS1
2. MS1 valida formato y guarda en MongoDB
3. MS1 publica mensaje a RabbitMQ
4. MS2 consume mensaje y valida cálculos
5. MS2 obtiene clave privada desde Vault
6. MS2 genera firma RSA y actualiza documento

## Requisitos

- Docker
- Docker Compose

## Levantar el Sistema

```bash
docker-compose up --build
```

## Portales de Acceso

Una vez levantado el sistema, puedes acceder a los siguientes portales:

### API y Documentación
- API REST MS1: http://localhost:5000
- Swagger UI (MS1): http://localhost:5000/swagger/index.html

### Gestión de Infraestructura
- RabbitMQ Management: http://localhost:15672
  - Usuario: admin
  - Password: admin123
  - Gestión de colas, exchanges y mensajes

- Vault UI: http://localhost:8200
  - Token: efact-root-token-dev-only
  - Gestión de secrets y claves RSA

### Monitoreo y Logs
- Elasticsearch: http://localhost:9200
  - API REST para consultas directas
  - Health: http://localhost:9200/_cluster/health

- Kibana: http://localhost:5601
  - Visualización de logs y métricas
  - Dashboards y análisis de datos

## Validaciones

MS1 valida formato:
- ID Documento: ABCD-012345678 (4 letras + guion + 9 dígitos)
- RUC: 11 dígitos numéricos
- Fecha: ISO 8601
- Previene duplicados

MS2 valida cálculos:
- IGV por item: precioTotal × 0.18
- IGV total: montoTotalSinImpuestos × 0.18
- Monto total: montoTotalSinImpuestos + igvTotal
- Tolerancia: 0.01

## Ejemplo de Uso

```bash
# Crear documento
curl -X POST http://localhost:5000/documents \
  -H "Content-Type: application/json" \
  -d '{
    "idDocumento": "ABCD-012345678",
    "rucEmisor": "20123456789",
    "rucReceptor": "20987654321",
    "fechaEmision": "2026-02-09T10:30:00Z",
    "montoTotalSinImpuestos": 1000.00,
    "igvTotal": 180.00,
    "montoTotal": 1180.00,
    "items": [
      {
        "descripcion": "Producto A",
        "precioUnitario": 100.00,
        "cantidad": 10,
        "precioTotal": 1000.00,
        "igvTotal": 180.00
      }
    ]
  }'

# Listar documentos
curl http://localhost:5000/documents

# Obtener documento específico
curl http://localhost:5000/documents/ABCD-012345678
```

## Detener

```bash
docker-compose down
```

Para eliminar volúmenes:

```bash
docker-compose down -v
```
