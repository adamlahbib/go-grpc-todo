package grpc

import (
	"context"
	"errors"

	pb "github.com/adamlahbib/go-grpc-todo/api/proto/v1"
	"github.com/adamlahbib/go-grpc-todo/internal/models"
	interfaces "github.com/adamlahbib/go-grpc-todo/pkg/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoServiceServer struct {
	useCase interfaces.UseCaseInterface
	pb.UnimplementedToDoServiceServer
}

func NewTodoServiceServer(useCase interfaces.UseCaseInterface) *TodoServiceServer {
	return &TodoServiceServer{useCase: useCase}
}

func (s *TodoServiceServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	// convert the deadline to a go time.Time
	deadline := req.GetTodo().GetDeadline() // Assuming deadline is of type *timestamppb.Timestamp
	if deadline == nil {
		return nil, status.Error(codes.InvalidArgument, "deadline field is missing")
	}

	// Validate and convert the timestamp
	if err := deadline.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "deadline field has invalid format: "+err.Error())
	}

	// convert the request to a todo model
	todo := models.Todo{
		Title:       req.GetTodo().Title,
		Description: req.GetTodo().Description,
		Deadline:    deadline.AsTime(),
	}

	if todo.Title == "" {
		return &pb.CreateResponse{}, errors.New("title is required")
	}

	// create the todo using the usecase
	createdTodo, err := s.useCase.Create(todo)
	if err != nil {
		return &pb.CreateResponse{}, err
	}

	// convert the created todo to a CreateResponse
	return &pb.CreateResponse{
		Id: int64(createdTodo.ID),
	}, nil
}
