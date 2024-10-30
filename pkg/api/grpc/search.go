package grpcApi

import (
	"context"

	"github.com/damione1/GoFinder/pkg/pb"
)

func (server *Server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	response := &pb.SearchResponse{
		Responses: []*pb.Response{
			{
				Name:  "Result 1",
				Type:  "Type 1",
				Score: 0.9,
			},
			{
				Name:  "Result 2",
				Type:  "Type 2",
				Score: 0.3,
			},
		},
	}

	return response, nil
}
