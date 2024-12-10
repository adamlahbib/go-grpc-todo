package usecase

import (
	"errors"

	"github.com/adamlahbib/go-grpc-todo/internal/models"
	interfaces "github.com/adamlahbib/go-grpc-todo/pkg/v1"
	"gorm.io/gorm"
)

type UseCase struct {
	repo interfaces.RepoInterface
}

func New(repo interfaces.RepoInterface) interfaces.UseCaseInterface {
	return &UseCase{repo}
}

func (u *UseCase) Create(todo models.Todo) (models.Todo, error) {
	return u.repo.Create(todo)
}

func (u *UseCase) Get(id int) (models.Todo, error) {
	var todo models.Todo
	var err error

	if todo, err = u.repo.Get(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Todo{}, errors.New("TODO record not found")
		}
		return models.Todo{}, err
	}
	return todo, nil
}

func (u *UseCase) Update(updatedTodo models.Todo) error {
	// check if todo exists
	if _, err := u.repo.Get(updatedTodo.Id); err != nil {
		return err
	}

	return u.repo.Update(updatedTodo)
}

func (u *UseCase) Delete(id int) error {
	// check if todo exists
	if _, err := u.repo.Get(id); err != nil {
		return err
	}
	return u.repo.Delete(id)
}

func (u *UseCase) GetAll() ([]models.Todo, error) {
	return u.repo.GetAll()
}
