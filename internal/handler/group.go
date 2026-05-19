package handler

import (
	"encoding/json"
	"net/http"

	"github.com/v3lichko/student-distribution/internal/models"
	"github.com/v3lichko/student-distribution/internal/response"
	"github.com/v3lichko/student-distribution/internal/storage"
)

type GroupHandler struct {
	storage *storage.GroupStorage
}

func NewGroupHandler(groupStorage *storage.GroupStorage) *GroupHandler {
	return &GroupHandler{storage: groupStorage}
}

// @Summary Create group
// @Tags groups
// @Accept json
// @Produce json
// @Param body body models.Group true "Group data"
// @Success 201 {object} models.Group
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /groups [post]
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if err := h.storage.CreateGroup(&group); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, group)
}

// @Summary Get all groups
// @Tags groups
// @Produce json
// @Success 200 {array} models.Group
// @Failure 500 {object} map[string]string
// @Router /groups [get]
func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.storage.GetGroups()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		return
	}

	response.WriteJSON(w, http.StatusOK, groups)
}
