package grpcApi

import (
	"fmt"

	"github.com/damione1/GoFinder/pkg/pb"
	"github.com/damione1/GoFinder/pkg/storage"
	"github.com/damione1/GoFinder/pkg/util"
	"gocloud.dev/blob"
)

type Server struct {
	pb.UnimplementedSearchServiceServer
	config util.Config
	bucket *blob.Bucket
}

func NewServer(config util.Config) (*Server, error) {
	var err error
	server := &Server{
		config: config,
	}

	switch config.Environment {
	case "production":
		//To be implemented
	case "development":
		server.bucket, err = storage.NewMinioBlobStorage("minio:9000", "minio", "miniosecret", "local-bucket", false) //Default Minio credentials
		if err != nil {
			return nil, fmt.Errorf("failed to create minio storage. %v", err)
		}
	default:
		return nil, fmt.Errorf("unknown environment: %s", config.Environment)
	}

	return server, nil
}

func (s *Server) Close() error {
	if s.bucket != nil {
		return s.bucket.Close()
	}
	return nil
}
