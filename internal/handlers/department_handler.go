package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

type UpdateDepartmentRequest struct {
	Name     *string         `json:"name"`
	ParentID json.RawMessage `json:"parent_id"`
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

	case errors.Is(err, services.ErrDepartmentNotFound):
		WriteError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, services.ErrDepartmentCannotBeParentOfItself):
		WriteError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, services.ErrDepartmentCycleDetected):
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

func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	depth, err := parseDepth(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	includeEmployees, err := parseIncludeEmployees(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.departmentService.GetDepartmentTree(id, depth, includeEmployees)
	if err != nil {
		h.handleDepartmentError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, response)
}

func parseDepth(r *http.Request) (int, error) {
	depthValue := r.URL.Query().Get("depth")
	if depthValue == "" {
		return 1, nil
	}

	depth, err := strconv.Atoi(depthValue)
	if err != nil {
		return 0, errors.New("depth must be an integer")
	}

	if depth < 0 {
		return 0, errors.New("depth must be greater than or equal to 0")
	}

	if depth > 5 {
		return 0, errors.New("depth must be less than or equal to 5")
	}

	return depth, nil
}

func parseIncludeEmployees(r *http.Request) (bool, error) {
	value := r.URL.Query().Get("include_employees")
	if value == "" {
		return true, nil
	}

	includeEmployees, err := strconv.ParseBool(value)
	if err != nil {
		return false, errors.New("include_employees must be true or false")
	}

	return includeEmployees, nil
}

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 {
		WriteError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	var request UpdateDepartmentRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	parentIDProvided := request.ParentID != nil
	var parentID *int

	if parentIDProvided {
		if bytes.Equal(request.ParentID, []byte("null")) {
			parentID = nil
		} else {
			var parsedParentID int

			if err := json.Unmarshal(request.ParentID, &parsedParentID); err != nil {
				WriteError(w, http.StatusBadRequest, "parent_id must be integer or null")
				return
			}

			if parsedParentID <= 0 {
				WriteError(w, http.StatusBadRequest, "parent_id must be greater than 0")
				return
			}

			parentID = &parsedParentID
		}
	}

	department, err := h.departmentService.UpdateDepartment(id, request.Name, parentIDProvided, parentID)
	if err != nil {
		h.handleDepartmentError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, department)
}
