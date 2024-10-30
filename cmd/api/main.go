package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	grpcApi "github.com/damione1/GoFinder/pkg/api/grpc"
	"github.com/damione1/GoFinder/pkg/pb"
	"github.com/damione1/GoFinder/pkg/util"
)

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("👋 Failed to load config")
	}

	go runGatewayServer(config)
	runGrpcServer(config)
}

func runGrpcServer(config util.Config) {
	log.Print("🍩 Starting gRPC server...")
	server, err := grpcApi.NewServer(config)
	if err != nil {
		log.Print(fmt.Sprintf("Failed to create gRPC server. %v", err))
	}
	defer server.Close()
	log.Print("🍩 gRPC server created")
	gprcLogger := grpc.UnaryInterceptor(grpcApi.GrpcLogger)
	grpcServer := grpc.NewServer(gprcLogger)
	pb.RegisterSearchServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Print("🍩 Starting to listen on port " + config.GRPCServerPort)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.GRPCServerPort))
	if err != nil {
		log.Print(fmt.Sprintf("🍩 Failed to listen. %v", err))
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Print(fmt.Sprintf("🍩 Failed to serve gRPC server over port %s. %v", listener.Addr().String(), err))
	}
}

func runGatewayServer(config util.Config) {
	log.Print("🍦 Starting HTTP server...")
	server, err := grpcApi.NewServer(config)
	if err != nil {
		log.Print(fmt.Sprintf("🍦 Failed to create HTTP server. %v", err))
	}
	defer server.Close()

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSearchServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("🍦 Failed to register HTTP gateway server.")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger", fs))
	log.Print(fmt.Sprintf("🍨 Swagger UI server started on http://localhost:%s/swagger/", config.HTTPServerPort))

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.HTTPServerPort))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen.")
	}

	log.Print(fmt.Sprintf("🍦 HTTP server started on http://localhost:%s/v1/", config.HTTPServerPort))
	handler := grpcApi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg(fmt.Sprintf("🍦 Failed to serve HTTP gateway server over port %s.", listener.Addr().String()))
	}
}
