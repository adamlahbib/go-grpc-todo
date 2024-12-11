package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/adamlahbib/go-grpc-todo/pkg/v1/handler"
	repo "github.com/adamlahbib/go-grpc-todo/pkg/v1/repository"
	"github.com/adamlahbib/go-grpc-todo/pkg/v1/usecase"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func InitiateServer(ctx context.Context, db *gorm.DB, port string) error {
	// create a listener
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	handler.NewTodoServiceServer(server, usecase.New(repo.New(db)))

	// graceful shutdown implementation
	c := make(chan os.Signal, 1)   // create a channel to listen for signals
	signal.Notify(c, os.Interrupt) // notify the channel when an interrupt signal is received

	go func() {
		for range c { // loop through the channel to receive signals
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
