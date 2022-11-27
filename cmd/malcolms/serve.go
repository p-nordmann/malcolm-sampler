package main

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	pb "github.com/p-nordmann/malcolm-sampler/grpc"
)

var port int
var batchSize int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start sampling service.",
	Long: `Starts grpc importance sampling service.
It expects to receive information about specific densities and offers to generate samples according
to a staircase approximation of said density.

The definition of the grpc service can be found in the package
github.com/p-nordmann/malcolm-sampler/grpc.

Example:
malcolms serve --port 1234`,

	Run: func(cmd *cobra.Command, args []string) {

		if batchSize < 1 {
			log.Fatalf("batch_size must be positive")
		}
		log.Printf("batch_size=%d", batchSize)

		log.Printf("listening on port %d...", port)
		listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// Options can be specified in call to NewServer.
		grpcServer := grpc.NewServer()

		pb.RegisterMalcolmSamplerServer(
			grpcServer,
			NewServer(batchSize),
		)
		grpcServer.Serve(listener)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 7352, "Port to listen to")
	serveCmd.Flags().IntVarP(&batchSize, "batch-size", "b", 1, "Size of batches returned by MakeSamples")
}
