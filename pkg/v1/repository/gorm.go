package repo

import (
	"github.com/adamlahbib/go-grpc-todo/internal/models"
	interfaces "github.com/adamlahbib/go-grpc-todo/pkg/v1"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) interfaces.RepoInterface {
	return &Repo{db}
}

func (r *Repo) Create(todo models.Todo) (models.Todo, error) {
	return todo, r.db.Create(&todo).Error
}

func (r *Repo) Get(id int) (models.Todo, error) {
	var todo models.Todo
	return todo, r.db.First(&todo, id).Error
}

func (r *Repo) Update(updatedTodo models.Todo) error {
	var todo models.Todo
	if err := r.db.First(&todo, updatedTodo.ID).Error; err != nil {
		return err // record not found
	}
	todo.Title = updatedTodo.Title
	todo.Description = updatedTodo.Description
	todo.Deadline = updatedTodo.Deadline
	return r.db.Save(&todo).Error
}

func (r *Repo) Delete(id int) error {
	return r.db.Delete(&models.Todo{}, id).Error
}

func (r *Repo) GetAll() ([]models.Todo, error) {
	var todos []models.Todo
	return todos, r.db.Find(&todos).Error
}
