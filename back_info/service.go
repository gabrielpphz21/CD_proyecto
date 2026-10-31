package main

import (
	"context"
	"fmt"
	"log"
	pb "protoc/proto_gen"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EnergyTransactionServer struct {
	pb.UnimplementedEnergyTransactionServiceServer
	conn *pgx.Conn
}

type TenantServer struct {
	pb.UnimplementedTenantServiceServer
	conn *pgx.Conn
}

type TakeOffServer struct {
	pb.UnimplementedTakeOffServiceServer
	conn *pgx.Conn
}

type DestinationServer struct {
	pb.UnimplementedDestinationServiceServer
	conn *pgx.Conn
}

// ---------- To_proto ----------

func (t *Tenant) To_proto() *pb.Tenant {
	return &pb.Tenant{
		Id:    int32(t.Id),
		Name:  t.Name,
		TaxId: t.Tax_id,
	}
}

func (et *EnergyTransaction) To_proto() *pb.EnergyTransaction {
	return &pb.EnergyTransaction{
		Id:          int32(et.ID),
		Volume:      float32(et.Volume),
		Origin:      int32(et.Origin),
		Destination: int32(et.Destination), // antes llamaba a et.TO_proto() (recursión infinita)
		State:       et.State,
		Date:        et.Date,
	}
}

func (d *Destination) To_proto() *pb.Destination {
	return &pb.Destination{
		Id:     int32(d.Id),
		Planet: d.Planet,
		State:  d.State,
	}
}

func (to *TakeOff) To_proto() *pb.TakeOff {
	return &pb.TakeOff{
		Id:     int32(to.Id),
		Rocket: int32(to.Rocket),
		Date:   to.Date,
		Runway: int32(to.Runway),
	}
}

// ---------- Tenant ----------

func (s *TenantServer) GetTenant(ctx context.Context, r *pb.GetTenantRequest) (*pb.GetTenantResponse, error) {
	tenants, err := get_Tenants(s.conn)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando tenants")
	}

	var tenants_pb []*pb.Tenant
	for _, t := range tenants {
		tenants_pb = append(tenants_pb, t.To_proto())
	}
	return &pb.GetTenantResponse{Tenants: tenants_pb}, nil
}

func (s *TenantServer) CreateTenant(ctx context.Context, r *pb.CreateTenantRequest) (*pb.CreateTenantResponse, error) {
	if r.Tenant == nil {
		return nil, status.Error(codes.InvalidArgument, "el tenant es requerido")
	}
	t := &Tenant{}
	t.Id = int(r.Tenant.Id)
	t.Name = r.Tenant.Name
	t.Tax_id = r.Tenant.TaxId
	if err := insert_tenant(t, s.conn); err != nil {
		return nil, status.Error(codes.Unavailable, "error creando tenant")
	}
	return &pb.CreateTenantResponse{Success: true}, nil
}

func (s *TenantServer) GetTenantId(ctx context.Context, r *pb.GetTenantIdRequest) (*pb.GetTenantIdResponse, error) {
	tnt, err := get_TenantId(s.conn, fmt.Sprint(r.Id))
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando tenants")
	}
	if tnt == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant %d no existe", r.Id)
	}

	return &pb.GetTenantIdResponse{Tenant: tnt.To_proto()}, nil
}

func (s *TenantServer) UpdateTenant(ctx context.Context, r *pb.UpdateTenantRequest) (*pb.UpdateTenantResponse, error) {
	if r.Tenant == nil {
		return nil, status.Error(codes.InvalidArgument, "el tenant es requerido")
	}
	t := &Tenant{}
	t.Id = int(r.Tenant.Id)
	t.Name = r.Tenant.Name
	t.Tax_id = r.Tenant.TaxId
	if err := update_tenant(t, s.conn); err != nil {
		log.Print("Error actualizando el tenant: ", err)
		return nil, status.Error(codes.Unavailable, "error actualizando tenant")
	}
	return &pb.UpdateTenantResponse{Success: true}, nil
}

func (s *TenantServer) DeleteTenant(ctx context.Context, r *pb.DeleteTenantRequest) (*pb.DeleteTenantResponse, error) {
	if err := delete_tenant(fmt.Sprint(r.Id), s.conn); err != nil {
		log.Print("error borrando tenant: ", err)
		return nil, status.Error(codes.Unavailable, "error borrando tenant")
	}
	return &pb.DeleteTenantResponse{Success: true}, nil
}

// ---------- EnergyTransaction ----------

func (s *EnergyTransactionServer) GetEnergyTransaction(ctx context.Context, r *pb.GetEnergyTransactionRequest) (*pb.GetEnergyTransactionResponse, error) {
	txs, err := get_energyTransactions(s.conn)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando transacciones de energía")
	}

	var txs_pb []*pb.EnergyTransaction
	for _, et := range txs {
		txs_pb = append(txs_pb, et.To_proto())
	}
	return &pb.GetEnergyTransactionResponse{EnergyTrasactions: txs_pb}, nil
}

func (s *EnergyTransactionServer) GetEnergyTransactionId(ctx context.Context, r *pb.GetEnergyTransactionIdRequest) (*pb.GetEnergyTransactionIdResponse, error) {
	et, err := get_energyTransactionsId(s.conn, fmt.Sprint(r.Id))
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando transacciones de energia")
	}
	if et == nil {
		return nil, status.Errorf(codes.NotFound, "EnergyTransaction %d no existe", r.Id)
	}
	return &pb.GetEnergyTransactionIdResponse{EnergyTransaction: et.To_proto()}, nil
}

func (s *EnergyTransactionServer) CreateEnergyTransaction(ctx context.Context, r *pb.CreateEnergyTransactionRequest) (*pb.CreateEnergyTransactionResponse, error) {
	if r.EnergyTransaction == nil {
		return nil, status.Error(codes.InvalidArgument, "la transacción de energía es requerida")
	}
	et := &EnergyTransaction{}
	et.ID = int(r.EnergyTransaction.Id)
	et.Volume = float64(r.EnergyTransaction.Volume)
	et.Origin = int(r.EnergyTransaction.Origin)
	et.Destination = int(r.EnergyTransaction.Destination)
	et.State = r.EnergyTransaction.State
	et.Date = r.EnergyTransaction.Date
	if err := insert_energyTransaction(et, s.conn); err != nil {
		return nil, status.Error(codes.Unavailable, "error creando transacción de energía")
	}
	return &pb.CreateEnergyTransactionResponse{Success: true}, nil
}

func (s *EnergyTransactionServer) UpdateEnergyTransaction(ctx context.Context, r *pb.UpdateEnergyTransactionRequest) (*pb.UpdateEnergyTransactionResponse, error) {
	if r.EnergyTransaction == nil {
		return nil, status.Error(codes.InvalidArgument, "la transacción de energía es requerida")
	}
	et := &EnergyTransaction{}
	et.ID = int(r.EnergyTransaction.Id)
	et.Volume = float64(r.EnergyTransaction.Volume)
	et.Origin = int(r.EnergyTransaction.Origin)
	et.Destination = int(r.EnergyTransaction.Destination)
	et.State = r.EnergyTransaction.State
	et.Date = r.EnergyTransaction.Date
	if err := update_energyTransaction(et, s.conn); err != nil {
		log.Print("error actualizando transacción de energía: ", err)
		return nil, status.Error(codes.Unavailable, "error actualizando transacción de energía")
	}
	return &pb.UpdateEnergyTransactionResponse{Success: true}, nil
}

func (s *EnergyTransactionServer) DeleteEnergyTransaction(ctx context.Context, r *pb.DeleteEnergyTransactionRequest) (*pb.DeleteEnergyTransactionResponse, error) {
	if err := delete_energyTransaction(fmt.Sprint(r.Id), s.conn); err != nil {
		log.Print("error borrando transacción de energía: ", err)
		return nil, status.Error(codes.Unavailable, "error borrando transacción de energía")
	}
	return &pb.DeleteEnergyTransactionResponse{Success: true}, nil
}

// ---------- Destination ----------

func (s *DestinationServer) GetDestination(ctx context.Context, r *pb.GetDestinationRequest) (*pb.GetDestinationResponse, error) {
	dests, err := get_destinations(s.conn)
	if err != nil {

		return nil, status.Error(codes.Unavailable, "error llamando destinos")
	}

	var dests_pb []*pb.Destination
	for _, d := range dests {
		dests_pb = append(dests_pb, d.To_proto())
	}
	return &pb.GetDestinationResponse{Destinations: dests_pb}, nil
}

func (s *DestinationServer) GetDestinationId(ctx context.Context, r *pb.GetDestinationIdRequest) (*pb.GetDestinationIdResponse, error) {
	dest, err := get_destinationsId(s.conn, fmt.Sprint(r.Id))
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando destinos")
	}
	if dest == nil {
		return nil, status.Errorf(codes.NotFound, "Destination %d no existe", r.Id)
	}

	return &pb.GetDestinationIdResponse{Destination: dest.To_proto()}, nil
}

func (s *DestinationServer) CreateDestination(ctx context.Context, r *pb.CreateDestinationRequest) (*pb.CreateDestinationResponse, error) {
	if r.Destination == nil {
		return nil, status.Error(codes.InvalidArgument, "el destino es requerido")
	}
	d := &Destination{}
	d.Id = int(r.Destination.Id)
	d.Planet = r.Destination.Planet
	d.State = r.Destination.State
	if err := insert_destination(d, s.conn); err != nil {
		return nil, status.Error(codes.Unavailable, "error creando destino")
	}
	return &pb.CreateDestinationResponse{Success: true}, nil
}

func (s *DestinationServer) UpdateDestination(ctx context.Context, r *pb.UpdateDestinationRequest) (*pb.UpdateDestinationResponse, error) {
	if r.Destination == nil {
		return nil, status.Error(codes.InvalidArgument, "el destino es requerido")
	}
	d := &Destination{}
	d.Id = int(r.Destination.Id)
	d.Planet = r.Destination.Planet
	d.State = r.Destination.State
	if err := update_destinantion(d, s.conn); err != nil {
		return nil, status.Error(codes.Unavailable, "error actualizando destino")
	}
	return &pb.UpdateDestinationResponse{Success: true}, nil
}

func (s *DestinationServer) DeleteDestination(ctx context.Context, r *pb.DeleteDestinationRequest) (*pb.DeleteDestinationResponse, error) {
	if err := delete_destination(fmt.Sprint(r.Id), s.conn); err != nil {
		return nil, status.Error(codes.Unavailable, "error borrando destino")
	}
	return &pb.DeleteDestinationResponse{Success: true}, nil
}

// ---------- TakeOff ----------

func (s *TakeOffServer) GetTakeOff(ctx context.Context, r *pb.GetTakeOffRequest) (*pb.GetTakeOffResponse, error) {
	takeoffs, err := get_takeoffs(s.conn)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando despegues")
	}

	var takeoffs_pb []*pb.TakeOff
	for _, to := range takeoffs {
		takeoffs_pb = append(takeoffs_pb, to.To_proto())
	}
	return &pb.GetTakeOffResponse{TakeOffs: takeoffs_pb}, nil
}

func (s *TakeOffServer) GetTakeOffId(ctx context.Context, r *pb.GetTakeOffIdRequest) (*pb.GetTakeOffIdResponse, error) {
	t_o, err := get_takeoffId(s.conn, fmt.Sprint(r.Id))
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando despegues")
	}
	if t_o == nil {
		return nil, status.Errorf(codes.NotFound, "EnergyTransaction %d no existe", r.Id)
	}

	return &pb.GetTakeOffIdResponse{TakeOff: t_o.To_proto()}, nil
}

func (s *TakeOffServer) CreateTakeOff(ctx context.Context, r *pb.CreateTakeOffRequest) (*pb.CreateTakeOffResponse, error) {
	if r.TakeOff == nil {
		return nil, status.Error(codes.InvalidArgument, "el despegue es requerido")
	}
	to := &TakeOff{}
	to.Id = int(r.TakeOff.Id)
	to.Rocket = int(r.TakeOff.Rocket)
	to.Date = r.TakeOff.Date
	to.Runway = int(r.TakeOff.Runway)
	if err := insert_takeoff(to, s.conn); err != nil {
		return nil, status.Error(codes.Unavailable, "error creando despegue")
	}
	return &pb.CreateTakeOffResponse{Success: true}, nil
}

func (s *TakeOffServer) UpdateTakeOff(ctx context.Context, r *pb.UpdateTakeOffRequest) (*pb.UpdateTakeOffResponse, error) {
	if r.TakeOff == nil {
		return nil, status.Error(codes.InvalidArgument, "el despegue es requerido")
	}
	to := &TakeOff{}
	to.Id = int(r.TakeOff.Id)
	log.Print("PROCESANDO TAKEOFF CON ID=: ", fmt.Sprint(r.TakeOff.Id))
	to.Rocket = int(r.TakeOff.Rocket)
	to.Date = r.TakeOff.Date
	to.Runway = int(r.TakeOff.Runway)
	if err := update_takeoff(to, s.conn); err != nil {
		log.Print("error actualizando despegue: ", err)
		log.Print("id: ", fmt.Sprint(to.Id))
		return nil, status.Error(codes.Unavailable, "error actualizando despegue")
	}
	return &pb.UpdateTakeOffResponse{Success: true}, nil
}

func (s *TakeOffServer) DeleteTakeOff(ctx context.Context, r *pb.DeleteTakeOffRequest) (*pb.DeleteTakeOffResponse, error) {
	if err := delete_takeoff(fmt.Sprint(r.Id), s.conn); err != nil {
		return nil, status.Error(codes.Unavailable, "error borrando despegue")
	}
	return &pb.DeleteTakeOffResponse{Success: true}, nil
}
