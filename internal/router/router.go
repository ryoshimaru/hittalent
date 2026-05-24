package router

import (
	"net/http"

	"github.com/ryoshimaru/hittalent/internal/handlers"
	"github.com/ryoshimaru/hittalent/internal/repositories"
	"github.com/ryoshimaru/hittalent/internal/services"

	"gorm.io/gorm"
)

func New(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()

	departmentRepository := repositories.NewDepartmentRepository(db)
	departmentService := services.NewDepartmentService(departmentRepository)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)

	employeeRepository := repositories.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeRepository(*employeeRepository, *departmentRepository)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("POST /departments/", departmentHandler.CreateDepartment)
	mux.HandleFunc("POST /departments/{id}/employees", employeeHandler.CreateEmployee)

	return mux
}
