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

type Server struct {
	pb.UnimplementedRocketServiceServer
	conn *pgx.Conn
}

func (r *Rocket) To_proto() *pb.Rocket {
	return &pb.Rocket{
		Id:       int32(r.Id),
		Satelite: int32(r.Satelite),
	}
}

func (hr *HyperRocket) To_proto() *pb.HyperRocket {
	return &pb.HyperRocket{
		Id:    int32(hr.Id),
		Fuel:  float32(hr.Fuel),
		Parts: hr.Parts,
		State: hr.State,
		X:     float32(hr.X),
		Y:     float32(hr.Y),
		Z:     float32(hr.Z),
		Date:  hr.Date,
	}
}

func (s *Server) GetRocket(ctx context.Context, r *pb.GetRocketRequest) (*pb.GetRocketResponse, error) {
	rockets, err := get_rockets(s.conn)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando cohetes")
	}

	var rockets_pb []*pb.Rocket
	for _, r := range rockets {
		rockets_pb = append(rockets_pb, r.To_proto())
	}
	return &pb.GetRocketResponse{
		Rockets: rockets_pb,
	}, nil

}

func (s *Server) GetRocketInfo(ctx context.Context, r *pb.GetRocketInfoRequest) (*pb.GetRocketInfoResponse, error) {
	rkt, hr, err := rocket_info(fmt.Sprint(r.Id), s.conn)

	if err != nil {
		return nil, status.Error(codes.Unavailable, "error llamando info de cohete")
	}
	var hr_pb []*pb.HyperRocket
	for _, hr_p := range hr {
		hr_pb = append(hr_pb, hr_p.To_proto())
	}
	return &pb.GetRocketInfoResponse{
		Rocket:      rkt.To_proto(),
		HyperRocket: hr_pb,
	}, nil
}

func (s *Server) CreateRocket(ctx context.Context, r *pb.CreateRocketRequest) (*pb.CreateRocketResponse, error) {

	rocket_p := &Rocket{}
	if r.Rocket == nil {
		return nil, status.Error(codes.InvalidArgument, "el cohete es requerido")
	}
	rocket_p.Id = int(r.Rocket.Id)
	rocket_p.Satelite = int(r.Rocket.Satelite)
	err1 := insert_rocket(rocket_p, s.conn)
	if err1 != nil {
		return nil, status.Error(codes.Unavailable, "error creando cohete")
	}
	return &pb.CreateRocketResponse{Success: true}, nil
}

func (s *Server) DeleteRocket(ctx context.Context, r *pb.DeleteRocketRequest) (*pb.DeleteRocketResponse, error) {

	err1 := delete_rocket(fmt.Sprint(r.Id), s.conn)
	if err1 != nil {
		return nil, status.Error(codes.Unavailable, "error borrando cohete")
	}
	return &pb.DeleteRocketResponse{Success: true}, nil
}

func (s *Server) UpdateRocket(ctx context.Context, r *pb.UpdateRocketRequest) (*pb.UpdateRocketResponse, error) {
	rocket_p := &Rocket{}
	if r.Rocket == nil {
		return nil, status.Error(codes.InvalidArgument, "el hyper cohete es requerido")
	}
	rocket_p.Id = int(r.Rocket.Id)
	rocket_p.Satelite = int(r.Rocket.Satelite)
	err1 := update_rocket(rocket_p, s.conn)
	if err1 != nil {
		return nil, status.Error(codes.Unavailable, "error actualizando cohete")
	}
	return &pb.UpdateRocketResponse{Success: true}, nil
}

func (s *Server) CreateHyperRocket(ctx context.Context, r *pb.CreateHyperRocketRequest) (*pb.CreateHyperRocketResponse, error) {
	hr_p := &HyperRocket{}
	if r.HyperRocket == nil {
		return nil, status.Error(codes.InvalidArgument, "el hyper cohete es requerido")

	}
	hr_p.Id = int(r.HyperRocket.Id)
	hr_p.Fuel = float64(r.HyperRocket.Fuel)
	hr_p.Date = r.HyperRocket.Date
	hr_p.Parts = r.HyperRocket.Parts
	hr_p.X = float64(r.HyperRocket.X)
	hr_p.Y = float64(r.HyperRocket.Y)
	hr_p.Z = float64(r.HyperRocket.Z)
	err1 := insert_hyperrocket(hr_p, s.conn)
	if err1 != nil {
		log.Print("error creando hyper cohete: ", err1)
		return nil, status.Error(codes.Unavailable, "error creando hyper cohete")

	}
	return &pb.CreateHyperRocketResponse{Success: true}, nil
}

func (s *Server) DeleteHyperRocket(ctx context.Context, r *pb.DeleteHyperRocketRequest) (*pb.DeleteHyperRocketResponse, error) {

	err1 := delete_hyperrocketLatest(fmt.Sprint(r.Id), s.conn)
	if err1 != nil {
		return nil, status.Error(codes.Unavailable, "error borrando hyper cohete")
	}
	return &pb.DeleteHyperRocketResponse{Success: true}, nil
}
