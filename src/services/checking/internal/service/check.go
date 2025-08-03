package service

import (
	"context"

	pb "github.com/Fl0rencess720/Majula/src/idl/checking"
)

func (s *CheckingService) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {

	return &pb.CheckResponse{}, nil
}
