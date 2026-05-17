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
func (h *GroupHandler) Groups(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.CreateGroup(w, r)
		return
	}
	if r.Method == http.MethodGet {
		h.GetGroups(w, r)
		return
	}
	response.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
		"error": "method not allowed",
	})
}

// @Summary Create group
// @Tags groups
// @Accept json
// @Produce json
// @Param body body models.Group true "Group data"
// @Success 201 {object} models.Group
// @Router /groups [post]
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	json.NewDecoder(r.Body).Decode(&group)

	err := h.storage.CreateGroup(&group)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError,
			map[string]string{"error": err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusCreated, group)
}

// @Summary Get groups
// @Tags groups
// @Produce json
// @Success 200 {array} models.Group
// @Router /groups [get]
func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.storage.GetGroups()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError,
			map[string]string{"error": err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, groups)
}
