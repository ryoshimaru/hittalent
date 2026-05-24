package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ryoshimaru/hittalent/internal/services"
)

type DepartmentHandler struct {
	departmentService *services.DepartmentService
}

func NewDepartmentHandler(departmentService *services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
	}
}

type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

func (h *DepartmentHandler) handleDepartmentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrDepartmentNameRequired):
		WriteError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, services.ErrDepartmentNameTooLong):
		WriteError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, services.ErrParentDepartmentNotFound):
		WriteError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, services.ErrDepartmentNameAlreadyExists):
		WriteError(w, http.StatusConflict, err.Error())

	default:
		WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var request CreateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	department, err := h.departmentService.CreateDepartment(request.Name, request.ParentID)
	if err != nil {
		h.handleDepartmentError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, department)
}
