package services

import (
	"errors"
	"strings"
	"time"

	"github.com/ryoshimaru/hittalent/internal/models"
	"github.com/ryoshimaru/hittalent/internal/repositories"
)

var (
	ErrEmployeeDepartmentDoesntExists = errors.New("department doesnt exists for this employee")
	ErrEmployeeFullNameEmpty          = errors.New("full name is empty")
	ErrEmployeeFullNameTooLong        = errors.New("employee full name too long")
	ErrEmployeePositionEmpty          = errors.New("employee position is empty")
	ErrEmployeePositionTooLong        = errors.New("employee positin too long")
)

type EmployeeService struct {
	employeeRepo   *repositories.EmployeeRepository
	departmentRepo *repositories.DepartmentRepository
}

func NewEmployeeRepository(employeeRepo repositories.EmployeeRepository, departmentRepo repositories.DepartmentRepository) *EmployeeService {
	return &EmployeeService{
		employeeRepo:   &employeeRepo,
		departmentRepo: &departmentRepo,
	}
}

func (s *EmployeeService) CreateEmployee(departmentId int, fullName string, position string, hired_at *time.Time) (*models.Employee, error) {
	if len(fullName) > 200 {
		return nil, ErrEmployeeFullNameTooLong
	}

	if len(position) > 200 {
		return nil, ErrEmployeePositionTooLong
	}

	fullName = strings.TrimSpace(fullName)
	position = strings.TrimSpace(position)

	if fullName == "" {
		return nil, ErrEmployeeFullNameEmpty
	}

	if position == "" {
		return nil, ErrEmployeePositionEmpty
	}

	exists, err := s.departmentRepo.ExistsByID(departmentId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrEmployeeDepartmentDoesntExists
	}

	employeeToCreate := *&models.Employee{
		DepartmentID: departmentId,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hired_at,
	}

	if err := s.employeeRepo.Create(&employeeToCreate); err != nil {
		return nil, err
	}

	return &employeeToCreate, nil
}
