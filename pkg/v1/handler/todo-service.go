package handler

import (
	"context"
	"errors"

	pb "github.com/adamlahbib/go-grpc-todo/api/proto/v1"
	"github.com/adamlahbib/go-grpc-todo/internal/models"
	interfaces "github.com/adamlahbib/go-grpc-todo/pkg/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TodoServiceServer struct {
	useCase interfaces.UseCaseInterface
	pb.UnimplementedToDoServiceServer
}

func NewTodoServiceServer(useCase interfaces.UseCaseInterface) *TodoServiceServer {
	return &TodoServiceServer{useCase: useCase}
}

func (s TodoServiceServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
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

func (s TodoServiceServer) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	// check if the id is provided
	if req.GetId() == 0 {
		return &pb.ReadResponse{}, errors.New("id is required")
	}

	// read the todo using the usecase
	todo, err := s.useCase.Get(int(req.GetId()))
	if err != nil {
		return &pb.ReadResponse{}, err
	}

	// convert deadline from time.time to *timestamppb.Timestamp
	deadline := timestamppb.New(todo.Deadline)

	// convert the todo to a ReadResponse
	return &pb.ReadResponse{
		Todo: &pb.ToDo{
			Id:          int64(todo.ID),
			Title:       todo.Title,
			Description: todo.Description,
			Deadline:    deadline,
		},
	}, nil
}

func (s TodoServiceServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	if req.GetTodo().GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id field is missing")
	}

	deadline := req.GetTodo().GetDeadline()
	if deadline == nil {
		return nil, status.Error(codes.InvalidArgument, "deadline field is missing")
	}

	if err := deadline.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "deadline field has invalid format: "+err.Error())
	}

	// convert the request to a todo model
	todo := models.Todo{
		Id:          int(req.GetTodo().Id),
		Title:       req.GetTodo().Title,
		Description: req.GetTodo().Description,
		Deadline:    deadline.AsTime(),
	}

	// update the todo using the usecase
	if err := s.useCase.Update(todo); err != nil {
		return &pb.UpdateResponse{
			Updated: false,
		}, err
	}

	return &pb.UpdateResponse{
		Updated: true,
	}, nil
}

func (s TodoServiceServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id field is missing")
	}

	// delete the todo using the usecase
	if err := s.useCase.Delete(int(req.GetId())); err != nil {
		return &pb.DeleteResponse{
			Deleted: false,
		}, err
	}

	return &pb.DeleteResponse{
		Deleted: true,
	}, nil
}

func (s TodoServiceServer) ReadAll(ctx context.Context, req *pb.ReadAllRequest) (*pb.ReadAllResponse, error) {
	// read all todos using the usecase
	todos, err := s.useCase.GetAll()
	if err != nil {
		return &pb.ReadAllResponse{}, err
	}

	// convert the todos to a ReadAllResponse
	var pbTodos []*pb.ToDo
	for _, todo := range todos {
		pbTodos = append(pbTodos, &pb.ToDo{
			Id:          int64(todo.ID),
			Title:       todo.Title,
			Description: todo.Description,
			Deadline:    timestamppb.New(todo.Deadline),
		})
	}

	return &pb.ReadAllResponse{
		Todos: pbTodos,
	}, nil
}
