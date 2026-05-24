package service

import (
	"errors"
	"strings"

	"github.com/ryoshimaru/hittalent/internal/models"
	"github.com/ryoshimaru/hittalent/internal/repositories"
)

var (
	ErrDepartmentNameAlreadyExists = errors.New("department name already exists in this parent")
	ErrDepartmentNameRequired      = errors.New("department name is required")
	ErrDepartmentNameTooLong       = errors.New("department name is too long")
	ErrParentDepartmentNotFound    = errors.New("parent department not found")
)

type DepartmentService struct {
	departmentRepo *repositories.DepartmentRepository
}

func NewDepartmentService(departmentRepo *repositories.DepartmentRepository) *DepartmentService {
	return &DepartmentService{
		departmentRepo: departmentRepo,
	}
}

func (d *DepartmentService) CreateDepartment(name string, parentID *int) (*models.Department, error) {
	if len(name) > 200 {
		return nil, ErrDepartmentNameTooLong
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrDepartmentNameRequired
	}

	if parentID != nil {
		exists, err := d.departmentRepo.ExistsByID(*parentID)

		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, ErrParentDepartmentNotFound
		}
	}

	nameExists, err := d.departmentRepo.ExistsNameInParent(name, parentID)
	if err != nil {
		return nil, err
	}

	if nameExists {
		return nil, ErrDepartmentNameAlreadyExists
	}

	departmentToCreate := &models.Department{
		Name:     name,
		ParentID: parentID,
	}

	if err := d.departmentRepo.Create(departmentToCreate); err != nil {
		return nil, err
	}

	return departmentToCreate, nil
}
