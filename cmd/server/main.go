package main

import (
	"context"
	"log"

	dbConfig "github.com/adamlahbib/go-grpc-todo/internal/db"
	"github.com/adamlahbib/go-grpc-todo/internal/models"
	server "github.com/adamlahbib/go-grpc-todo/pkg/protocol/grpc"
)

func main() {

	// connect to the database
	db := dbConfig.DbConn()

	// migrate the database
	if err := db.AutoMigrate(&models.Todo{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := server.InitiateServer(
		context.Background(),
		db,
		"8080",
	); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

}
