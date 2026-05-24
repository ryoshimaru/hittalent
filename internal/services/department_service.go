package services

import (
	"errors"
	"strings"

	"github.com/ryoshimaru/hittalent/internal/models"
	"github.com/ryoshimaru/hittalent/internal/repositories"
)

var (
	ErrDepartmentCannotBeParentOfItself = errors.New("department cannot be parent of itself")
	ErrDepartmentCycleDetected          = errors.New("department cycle detected")

	ErrDepartmentNotFound          = errors.New("department not found")
	ErrDepartmentNameAlreadyExists = errors.New("department name already exists in this parent")
	ErrDepartmentNameRequired      = errors.New("department name is required")
	ErrDepartmentNameTooLong       = errors.New("department name is too long")
	ErrParentDepartmentNotFound    = errors.New("parent department not found")
)

type DepartmentService struct {
	departmentRepo *repositories.DepartmentRepository
	employeeRepo   *repositories.EmployeeRepository
}

func NewDepartmentService(departmentRepo *repositories.DepartmentRepository, employeeRepo *repositories.EmployeeRepository) *DepartmentService {
	return &DepartmentService{
		departmentRepo: departmentRepo,
		employeeRepo:   employeeRepo,
	}
}

type DepartmentTreeResponse struct {
	Department models.Department        `json:"department"`
	Employees  *[]models.Employee       `json:"employees,omitempty"`
	Children   []DepartmentTreeResponse `json:"children"`
}

func (s *DepartmentService) buildDepartmentTree(department models.Department, depth int, includeEmployees bool) (*DepartmentTreeResponse, error) {
	response := &DepartmentTreeResponse{
		Department: department,
		Children:   make([]DepartmentTreeResponse, 0),
	}

	if includeEmployees {
		employees, err := s.employeeRepo.GetByDepartmentID(department.ID)
		if err != nil {
			return nil, err
		}

		response.Employees = &employees
	}

	if depth <= 0 {
		return response, nil
	}

	children, err := s.departmentRepo.GetChildrenByParentID(department.ID)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		childResponse, err := s.buildDepartmentTree(child, depth-1, includeEmployees)
		if err != nil {
			return nil, err
		}

		response.Children = append(response.Children, *childResponse)
	}

	return response, nil
}

func (s *DepartmentService) GetDepartmentTree(id int, depth int, includeEmployees bool) (*DepartmentTreeResponse, error) {
	department, err := s.departmentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if department == nil {
		return nil, ErrDepartmentNotFound
	}

	return s.buildDepartmentTree(*department, depth, includeEmployees)
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

func (s *DepartmentService) wouldCreateCycle(departmentID int, newParentID int) (bool, error) {
	currentParentID := &newParentID

	for currentParentID != nil {
		if *currentParentID == departmentID {
			return true, nil
		}

		parent, err := s.departmentRepo.GetByID(*currentParentID)
		if err != nil {
			return false, err
		}

		if parent == nil {
			return false, ErrParentDepartmentNotFound
		}

		currentParentID = parent.ParentID
	}

	return false, nil
}

func (s *DepartmentService) UpdateDepartment(id int, name *string, parentIDProvided bool, parentID *int) (*models.Department, error) {
	department, err := s.departmentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if department == nil {
		return nil, ErrDepartmentNotFound
	}

	newName := department.Name
	newParentID := department.ParentID

	if name != nil {
		trimmedName := strings.TrimSpace(*name)

		if trimmedName == "" {
			return nil, ErrDepartmentNameRequired
		}

		if len(trimmedName) > 200 {
			return nil, ErrDepartmentNameTooLong
		}

		newName = trimmedName
	}

	if parentIDProvided {
		if parentID != nil {
			if *parentID == id {
				return nil, ErrDepartmentCannotBeParentOfItself
			}

			parentExists, err := s.departmentRepo.ExistsByID(*parentID)
			if err != nil {
				return nil, err
			}

			if !parentExists {
				return nil, ErrParentDepartmentNotFound
			}

			hasCycle, err := s.wouldCreateCycle(id, *parentID)
			if err != nil {
				return nil, err
			}

			if hasCycle {
				return nil, ErrDepartmentCycleDetected
			}
		}

		newParentID = parentID
	}

	nameExists, err := s.departmentRepo.ExistsByNameAndParentIDExceptID(newName, newParentID, id)
	if err != nil {
		return nil, err
	}

	if nameExists {
		return nil, ErrDepartmentNameAlreadyExists
	}

	department.Name = newName
	department.ParentID = newParentID

	if err := s.departmentRepo.Update(department); err != nil {
		return nil, err
	}

	return department, nil
}
