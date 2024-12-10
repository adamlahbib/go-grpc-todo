package v1

import "github.com/adamlahbib/go-grpc-todo/internal/models"

type RepoInterface interface {
	Create(models.Todo) (models.Todo, error)
	Get(id int) (models.Todo, error)
	Update(models.Todo) error
	Delete(id int) error
	GetAll() ([]models.Todo, error)
}

type UseCaseInterface interface {
	Create(models.Todo) (models.Todo, error)
	Get(id int) (models.Todo, error)
	Update(models.Todo) error
	Delete(id int) error
	GetAll() ([]models.Todo, error)
}
