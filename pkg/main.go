package pkg

import (
	"fmt"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, World!")

	// This is just to make sure that the grpc package is used
	_ = grpc.NewServer()
}
