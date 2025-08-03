package service

import (
	"context"

	pb "github.com/Fl0rencess720/Majula/src/idl/checking"
)

func (s *CheckingService) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckListResponse, error) {
	res, err := s.cu.Check(req.Text)
	if err != nil {
		return nil, err
	}
	resp := &pb.CheckListResponse{}
	for _, e := range res {
		resp.Results = append(resp.Results, &pb.CheckResponse{
			OriginalText: e.OriginalText,
			Result:       e.Result,
			Sources:      e.Sources,
			Reason:       e.Reason,
		})
	}
	return resp, nil
}
