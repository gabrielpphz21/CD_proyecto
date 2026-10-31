#!/bin/bash
# =============================================================================
# Pruebas de la API gRPC con grpcurl
# Copia y pega cada comando individualmente en tu terminal.
#
# Requisitos:
#   - grpcurl instalado:  brew install grpcurl   (o)   go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
#   - Archivo .proto en la ruta indicada abajo (ajusta PROTO_FILE)
#   - Servidor gRPC corriendo (ajusta HOST)
#
# Si tu servidor tiene reflection habilitado, puedes quitar "-import-path" y "-proto"
# de cada comando.
# =============================================================================

repetir flags (pega esto también)
G="grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE"

# Si prefieres no usar la variable G, cada comando de abajo ya la usa como: $G -d '{...}' $HOST servicio/Metodo
# (si abres otra terminal, vuelve a pegar las líneas de configuración)

# Nota: los mensajes usan los nombres JSON (camelCase) de los campos, pero grpcurl
# también acepta snake_case (tax_id, hyper_satelite, etc.).


# =============================================================================
# 0. EXPLORACIÓN (opcional)
# =============================================================================

# Listar servicios (solo si hay reflection; con -proto usa describe)
$G $HOST list

# Listar métodos de un servicio
$G $HOST list SateliteService

# Describir un mensaje
$G $HOST describe Satelite


# =============================================================================
# 1. TENANT  (créalo primero, Satelite depende de él)
# =============================================================================

# Crear tenant
$G -d '{"tenant": { "name": "Acme Corp", "tax_id": "ACM123456XYZ"}}' $HOST TenantService/CreateTenant

# Crear un segundo tenant
$G -d '{"tenant": { "name": "Orbital Inc", "tax_id": "ORB987654ABC"}}' $HOST TenantService/CreateTenant

# Listar todos los tenants
$G -d '{}' $HOST TenantService/GetTenant

# Obtener tenant por id
$G -d '{"id": 2}' $HOST TenantService/GetTenantId

# Actualizar tenant
$G -d '{"tenant": {"id": 2, "name": "Acme Corp Actualizado", "tax_id": "ACM123456XYZ"}}' $HOST TenantService/UpdateTenant

# Verificar actualización
$G -d '{"id": 2}' $HOST TenantService/GetTenantId


# =============================================================================
# 2. SATELITE
# =============================================================================

# Crear satélite (tenant = 1)
$G -d '{"satelite": { "tenant": 1}}' $HOST SateliteService/CreateSatelite

# Crear otro satélite
$G -d '{"satelite": {"tenant": 2}}' $HOST SateliteService/CreateSatelite

# Listar todos los satélites
$G -d '{}' $HOST SateliteService/GetSatelite

# Actualizar satélite
$G -d '{"satelite": {"id": 10, "tenant": 2}}' $HOST SateliteService/UpdateSatelite

# --- HyperSatelite ---

# Crear hyper satélite
$G -d '{"hyper_satelite": {"id": 1, "energy": 1500.5, "parts": "{\"panel_solar\": 2, \"antena\": 1, \"bateria\": {\"tipo\": \"litio\", \"capacidad_wh\": 5000}}", "state": "active", "x": 10.5, "y": 20.3, "z": 400.0, "date": "2026-10-30"}}' $HOST SateliteService/CreateHyperSatelite
# Obtener info del satélite (incluye sus hyper satélites)
$G -d '{"id": 1}' $HOST SateliteService/GetSateliteInfo


# =============================================================================
# 3. ROCKET
# =============================================================================

# Crear cohete (satelite = 1)
$G -d '{"rocket": {"satelite": 2}}' $HOST RocketService/CreateRocket

# Crear otro cohete
$G -d '{"rocket": {"id": 2, "satelite": 2}}' $HOST RocketService/CreateRocket

# Listar todos los cohetes
$G -d '{}' $HOST RocketService/GetRocket

# Actualizar cohete
$G -d '{"rocket": {"id": 1, "satelite": 2}}' $HOST RocketService/UpdateRocket

# --- HyperRocket ---

# Crear hyper cohete
$G -d '{"hyper_rocket": {"id": 1, "fuel": 5500.5, "parts": "{\"motor\": {\"modelo\": \"M1\", \"empuje_kn\": 850}, \"tanque\": 1}", "state": "launched", "x": 5.0, "y": 5.0, "z": 100.0, "date": "2026-10-31"}}' $HOST RocketService/CreateHyperRocket
# Crear otro hyper cohete
$G -d '{"hyper_rocket": {"id": 2, "fuel": 5500.5, "parts": "motor,tanque", "state": "launched", "x": 5.0, "y": 5.0, "z": 100.0, "date": "2026-10-31"}}' $HOST RocketService/CreateHyperRocket

# Obtener info del cohete (incluye sus hyper cohetes)
$G -d '{"id": 1}' $HOST RocketService/GetRocketInfo


# =============================================================================
# 4. DESTINATION
# =============================================================================

# Crear destino
$G -d '{"destination": { "planet": "Marte", "state": "available"}}' $HOST DestinationService/CreateDestination

# Crear otro destino
$G -d '{"destination": {"planet": "Luna", "state": "reserved"}}' $HOST DestinationService/CreateDestination

# Listar todos los destinos
$G -d '{}' $HOST DestinationService/GetDestination

# Obtener destino por id
$G -d '{"id": 1}' $HOST DestinationService/GetDestinationId

# Actualizar destino
$G -d '{"destination": {"id": 1, "planet": "Marte", "state": "reserved"}}' $HOST DestinationService/UpdateDestination


# =============================================================================
# 5. TAKEOFF
# =============================================================================

# Crear despegue (rocket = 1)
$G -d '{"take_off": {"id": 1, "rocket": 1, "date": "2026-11-01", "runway": 3}}' $HOST TakeOffService/CreateTakeOff

# Crear otro despegue
$G -d '{"take_off": { "rocket": 1 , "date": "2026-11-05", "runway": 1}}' $HOST TakeOffService/CreateTakeOff

# Listar todos los despegues
$G -d '{}' $HOST TakeOffService/GetTakeOff

# Obtener despegue por id
$G -d '{"id": 1}' $HOST TakeOffService/GetTakeOffId

# Actualizar despegue
$G -d '{"take_off": {"id": 1, "rocket": 1, "date": "2026-11-02", "runway": 1}}' $HOST TakeOffService/UpdateTakeOff


# =============================================================================
# 6. ENERGY TRANSACTION  (origin / destination = ids de satélites)
# =============================================================================


# Crear otra transacción
$G -d '{"energy_transaction": {"volume": 100.0, "origin": 2, "destination": 1, "state": "completed", "date": "2026-10-31"}}' $HOST EnergyTransactionService/CreateEnergyTransaction

# Listar todas las transacciones
$G -d '{}' $HOST EnergyTransactionService/GetEnergyTransaction

# Obtener transacción por id
$G -d '{"id": 7}' $HOST EnergyTransactionService/GetEnergyTransactionId

# Actualizar transacción
$G -d '{"energy_transaction": {"id": 7, "volume": 250.75, "origin": 1, "destination": 2, "state": "completed", "date": "2026-10-30"}}' $HOST EnergyTransactionService/UpdateEnergyTransaction


# =============================================================================
# 7. DELETES  (en orden inverso a las dependencias)
# =============================================================================

# Borrar transacciones de energía
$G -d '{"id": 1}' $HOST EnergyTransactionService/DeleteEnergyTransaction
$G -d '{"id": 2}' $HOST EnergyTransactionService/DeleteEnergyTransaction

# Borrar despegues
$G -d '{"id": 1}' $HOST TakeOffService/DeleteTakeOff
$G -d '{"id": 2}' $HOST TakeOffService/DeleteTakeOff

# Borrar destinos
$G -d '{"id": 1}' $HOST DestinationService/DeleteDestination
$G -d '{"id": 2}' $HOST DestinationService/DeleteDestination

# Borrar hyper cohetes
$G -d '{"id": 1}' $HOST RocketService/DeleteHyperRocket
$G -d '{"id": 2}' $HOST RocketService/DeleteHyperRocket

# Borrar cohetes
$G -d '{"id": 1}' $HOST RocketService/DeleteRocket
$G -d '{"id": 2}' $HOST RocketService/DeleteRocket

# Borrar hyper satélites
$G -d '{"id": 1}' $HOST SateliteService/DeleteHyperSatelite
$G -d '{"id": 2}' $HOST SateliteService/DeleteHyperSatelite

# Borrar satélites
$G -d '{"id": 1}' $HOST SateliteService/DeleteSatelite
$G -d '{"id": 2}' $HOST SateliteService/DeleteSatelite

# Borrar tenants
$G -d '{"id": 1}' $HOST TenantService/DeleteTenant
$G -d '{"id": 2}' $HOST TenantService/DeleteTenant