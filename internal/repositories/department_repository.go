package repositories

import (
	"errors"

	"github.com/ryoshimaru/hittalent/internal/models"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{
		db: db,
	}
}

func (r *DepartmentRepository) Create(department *models.Department) error {
	return r.db.Create(department).Error
}

func (r *DepartmentRepository) ExistsByID(id int) (bool, error) {
	var department models.Department

	err := r.db.First(&department, id).Error
	if err == nil {
		return true, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

func (r *DepartmentRepository) ExistsNameInParent(name string, parentID *int) (bool, error) {
	var count int64

	query := r.db.Model(models.Department{}).Where("name = ?", name)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
