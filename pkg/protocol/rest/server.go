package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	pb "github.com/adamlahbib/go-grpc-todo/pkg/api/v1"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitiateServer(ctx context.Context, grpcPort, httpPort string) error {
	// create a new context with a cancel function to stop the server
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// create a new gRPC gateway mux server to serve the REST API
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// register the REST API handler for the ToDoService gRPC service on the mux server at the specified gRPC port
	if err := pb.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, "localhost"+grpcPort, opts); err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	// create a new HTTP server to serve the REST API
	srv := &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	// graceful shutdown of the HTTP server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down the HTTP server")
			_, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			_ = srv.Shutdown(ctx)
		}
	}()

	log.Println("starting HTTP/REST server gateway on port " + httpPort)
	return srv.ListenAndServe()
}
