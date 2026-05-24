package repositories

import (
	"github.com/ryoshimaru/hittalent/internal/models"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

func (r *EmployeeRepository) GetByDepartmentID(departmentID int) ([]models.Employee, error) {
	var employees []models.Employee

	err := r.db.
		Where("department_id = ?", departmentID).
		Order("full_name ASC").
		Find(&employees).
		Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}
