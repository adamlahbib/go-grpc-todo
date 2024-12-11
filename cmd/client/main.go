package main

import (
	"context"
	"log"
	"time"

	pb "github.com/adamlahbib/go-grpc-todo/api/proto/v1"
	lorem "github.com/derektata/lorem/ipsum"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// establish connection to the server
	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// create a new todo client
	client := pb.NewToDoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // set a timeout of 5 seconds for the request
	defer cancel()

	// create a new todo
	_time := time.Now().In(time.UTC)
	_deadline := timestamppb.New(_time)
	// random lorem ipsum text for the description and title
	g := lorem.NewGenerator()
	_title := g.Generate(10)
	_description := g.GenerateParagraphs(1)

	req1 := &pb.CreateRequest{
		Todo: &pb.ToDo{
			Title:       _title,
			Description: _description,
			Deadline:    _deadline,
		},
	}

	res1, err := client.Create(ctx, req1)
	if err != nil {
		log.Fatalf("failed to create todo: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)
	_id := res1.GetId()

	log.Printf("Created todo with id: %d\n\n", _id)

	// read the created todo
	req2 := &pb.ReadRequest{
		Id: _id,
	}
	res2, err := client.Read(ctx, req2)
	if err != nil {
		log.Fatalf("failed to read todo: %v", err)
	}
	log.Printf("Read result: <%+v>\n\n", res2)

	// update the created todo
	req3 := &pb.UpdateRequest{
		Todo: &pb.ToDo{
			Id:          res2.GetTodo().GetId(),
			Title:       res2.GetTodo().GetTitle() + " updated",
			Description: res2.GetTodo().GetDescription() + " updated",
			Deadline:    res2.GetTodo().GetDeadline(),
		},
	}
	res3, err := client.Update(ctx, req3)
	if err != nil {
		log.Fatalf("failed to update todo: %v", err)
	}
	log.Printf("Update result: <%+v>\n\n", res3)

	// call read all todos
	req4 := &pb.ReadAllRequest{}
	res4, err := client.ReadAll(ctx, req4)
	if err != nil {
		log.Fatalf("failed to read all todos: %v", err)
	}
	log.Printf("ReadAll result: <%+v>\n\n", res4)

	// delete the created todo
	req5 := &pb.DeleteRequest{
		Id: _id,
	}
	res5, err := client.Delete(ctx, req5)
	if err != nil {
		log.Fatalf("failed to delete todo: %v", err)
	}
	log.Printf("Delete result: <%+v>\n\n", res5)
}
