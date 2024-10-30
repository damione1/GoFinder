package grpcApi

import (
	"context"

	"github.com/damione1/GoFinder/pkg/pb"
)

func (server *Server) GetUploadURLRequest(ctx context.Context, req *pb.GetUploadURLRequest) (*pb.GetUploadURLResponse, error) {
	response := &pb.GetUploadURLResponse{
		UploadUrl: "http://localhost:8080/upload",
		FileId:    "1234",
	}

	return response, nil
}
