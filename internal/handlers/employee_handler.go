package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ryoshimaru/hittalent/internal/services"
)

type EmployeeHandler struct {
	employeeService *services.EmployeeService
}

func NewEmployeeHandler(employeeService *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

type CreateEmployeeRequest struct {
	FullName string  `json:"full_name"`
	Position string  `json:"position"`
	HiredAt  *string `json:"hired_at"`
}

func parseHiredAt(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	date := strings.TrimSpace(*value)
	if date == "" {
		return nil, errors.New("hired_at must be in YYYY-MM-DD format")
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("hired_at must be in YYYY-MM-DD format")
	}

	return &parsedDate, nil
}

func (h *EmployeeHandler) handleEmployeeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrEmployeeDepartmentDoesntExists):
		WriteError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, services.ErrEmployeeFullNameEmpty):
		WriteError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, services.ErrEmployeeFullNameTooLong):
		WriteError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, services.ErrEmployeePositionEmpty):
		WriteError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, services.ErrEmployeePositionTooLong):
		WriteError(w, http.StatusBadRequest, err.Error())

	default:
		WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	departmentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || departmentID < 0 {
		WriteError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	var request CreateEmployeeRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	hiredAt, err := parseHiredAt(request.HiredAt)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	employee, err := h.employeeService.CreateEmployee(departmentID, request.FullName, request.Position, hiredAt)
	if err != nil {
		h.handleEmployeeError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, employee)

}
