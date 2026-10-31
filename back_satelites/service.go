package main

import (
	"context"
	"fmt"
	pb "protoc/proto_gen"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedSateliteServiceServer
	conn *pgx.Conn
}

func (r *Satelite) To_proto() *pb.Satelite {
	return &pb.Satelite{
		Id:     int32(r.Id),
		Tenant: int32(r.Tenant),
	}
}

func (hr *HyperSatelite) To_proto() *pb.HyperSatelite {
	return &pb.HyperSatelite{
		Id:     int32(hr.Id),
		Energy: float32(hr.Energy),
		Parts:  hr.Parts,
		State:  hr.State,
		X:      float32(hr.X),
		Y:      float32(hr.Y),
		Z:      float32(hr.Z),
		Date:   hr.Date,
	}
}

func (s *Server) GetSatelite(ctx context.Context, r *pb.GetSateliteRequest) (*pb.GetSateliteResponse, error) {
	Satelites, err := get_satelites(s.conn)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando satelites")
	}

	var Satelites_pb []*pb.Satelite
	for _, r := range Satelites {
		Satelites_pb = append(Satelites_pb, r.To_proto())
	}
	return &pb.GetSateliteResponse{
		Satelites: Satelites_pb,
	}, nil

}

func (s *Server) GetSateliteInfo(ctx context.Context, r *pb.GetSateliteInfoRequest) (*pb.GetSateliteInfoResponse, error) {
	rkt, hr, err := satelite_info(fmt.Sprint(r.Id), s.conn)

	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando info de satelite")
	}
	var hr_pb []*pb.HyperSatelite
	for _, hr_p := range hr {
		hr_pb = append(hr_pb, hr_p.To_proto())
	}
	return &pb.GetSateliteInfoResponse{
		Satelite:      rkt.To_proto(),
		HyperSatelite: hr_pb,
	}, nil
}

func (s *Server) CreateSatelite(ctx context.Context, r *pb.CreateSateliteRequest) (*pb.CreateSateliteResponse, error) {
	Satelite_p := &Satelite{}
	if r.Satelite == nil {
		return nil, status.Error(codes.InvalidArgument, "el satelite es requerido")
	}
	Satelite_p.Id = int(r.Satelite.Id)
	Satelite_p.Tenant = int(r.Satelite.Tenant)
	err1 := insert_satelite(Satelite_p, s.conn)
	if err1 != nil {
		return nil, status.Error(codes.Unavailable, "error creando satelite")
	}
	return &pb.CreateSateliteResponse{Success: true}, nil
}

func (s *Server) DeleteSatelite(ctx context.Context, r *pb.DeleteSateliteRequest) (*pb.DeleteSateliteResponse, error) {

	err1 := delete_satelite(fmt.Sprint(r.Id), s.conn)
	if err1 != nil {
		return nil, status.Error(codes.Unavailable, "error borrando satelite")
	}
	return &pb.DeleteSateliteResponse{Success: true}, nil
}

func (s *Server) UpdateSatelite(ctx context.Context, r *pb.UpdateSateliteRequest) (*pb.UpdateSateliteResponse, error) {
	Satelite_p := &Satelite{}
	if r.Satelite == nil {
		return nil, status.Error(codes.InvalidArgument, "el satelite es requerido")
	}
	Satelite_p.Id = int(r.Satelite.Id)
	Satelite_p.Tenant = int(r.Satelite.Tenant)
	err1 := update_satelite(Satelite_p, s.conn)
	if err1 != nil {
		return nil, status.Error(codes.Unavailable, "error actualizando satelite")
	}
	return &pb.UpdateSateliteResponse{Success: true}, nil
}

func (s *Server) CreateHyperSatelite(ctx context.Context, r *pb.CreateHyperSateliteRequest) (*pb.CreateHyperSateliteResponse, error) {
	hr_p := &HyperSatelite{}
	if r.HyperSatelite == nil {
		return nil, status.Error(codes.InvalidArgument, "el hyper satelite es requerido")
	}
	hr_p.Id = int(r.HyperSatelite.Id)
	hr_p.Energy = int(r.HyperSatelite.Energy)
	hr_p.Date = r.HyperSatelite.Date
	hr_p.Parts = r.HyperSatelite.Parts
	hr_p.X = float64(r.HyperSatelite.X)
	hr_p.Y = float64(r.HyperSatelite.Y)
	hr_p.Z = float64(r.HyperSatelite.Z)
	err1 := insert_hypersatelites(hr_p, s.conn)
	if err1 != nil {
		return nil, status.Error(codes.Unavailable, "error creando hyper satelite")
	}
	return &pb.CreateHyperSateliteResponse{Success: true}, nil
}

func (s *Server) DeleteHyperSatelite(ctx context.Context, r *pb.DeleteHyperSateliteRequest) (*pb.DeleteHyperSateliteResponse, error) {

	err1 := delete_hypersateliteLatest(fmt.Sprint(r.Id), s.conn)
	if err1 != nil {
		return nil, status.Error(codes.Unavailable, "error borrando satelite")
	}
	return &pb.DeleteHyperSateliteResponse{Success: true}, nil
}
