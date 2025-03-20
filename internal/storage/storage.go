package storage

import "github.com/devarshkikani/demo_project/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudentData(id int64, name string, email string, age int) (int64, error)
}
