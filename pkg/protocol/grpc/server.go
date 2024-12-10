package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	pb "github.com/adamlahbib/go-grpc-todo/api/proto/v1"
	handler "github.com/adamlahbib/go-grpc-todo/pkg/v1/handler"
	"google.golang.org/grpc"
)

func InitiateServer(ctx context.Context, api handler.TodoServiceServer, port string) error {
	// create a listener
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	pb.RegisterToDoServiceServer(server, api)

	// graceful shutdown implementation
	c := make(chan os.Signal, 1)   // create a channel to listen for signals
	signal.Notify(c, os.Interrupt) // notify the channel when an interrupt signal is received

	go func() {
		for range c {
			// signal is ctrl+c, handle the shutdown
			log.Println("shutting down the gRPC server")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	// start the server
	log.Println("starting gRPC server on port " + port)
	return server.Serve(listen)
}
