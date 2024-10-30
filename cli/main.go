package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	"github.com/damione1/GoFinder/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var CLI struct {
	Find struct {
		Query string `arg:"" required:"" help:"Query to search for."`
	} `cmd:"" help:"Find."`
	Ingest struct {
		Name string `arg:"" required:"" help:"Name of the file to ingest."`
	} `cmd:"" help:"Ingest."`
}

func main() {
	ctx := kong.Parse(&CLI)

	// Set up a connection to the server.
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new SearchServiceClient
	grpcClient := pb.NewSearchServiceClient(conn)

	switch ctx.Command() {
	case "find <query>":
		// Call the Search method
		response, err := grpcClient.Search(context.Background(), &pb.SearchRequest{Query: &pb.Query{
			Search: CLI.Find.Query,
		}})
		if err != nil {
			log.Fatalf("Error calling Search: %v", err)
		}

		// Print the response
		fmt.Printf("Search results for query '%s':\n", CLI.Find.Query)
		for _, result := range response.GetResponses() {
			fmt.Printf("- %s\n", result)
		}
	case "ingest <name>":
		GetUploadURLResponse, err := grpcClient.GetUploadURL(context.Background(), &pb.GetUploadURLRequest{
			Name: CLI.Ingest.Name,
		})
		if err != nil {
			log.Fatalf("Error calling GetUploadURL: %v", err)
		}
		// Print the response

		fmt.Printf("Ingested file '%s'.\n", GetUploadURLResponse.UploadUrl)
	default:
		panic(ctx.Command())
	}
}
