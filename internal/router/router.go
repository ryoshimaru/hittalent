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
	employeeRepository := repositories.NewEmployeeRepository(db)
	departmentService := services.NewDepartmentService(db, departmentRepository, employeeRepository)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)

	employeeService := services.NewEmployeeService(*employeeRepository, *departmentRepository)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("POST /departments/", departmentHandler.CreateDepartment)
	mux.HandleFunc("POST /departments/{id}/employees", employeeHandler.CreateEmployee)
	mux.HandleFunc("GET /departments/{id}", departmentHandler.GetDepartment)
	mux.HandleFunc("PATCH /departments/{id}", departmentHandler.UpdateDepartment)
	mux.HandleFunc("DELETE /departments/{id}", departmentHandler.DeleteDepartment)

	return mux
}
