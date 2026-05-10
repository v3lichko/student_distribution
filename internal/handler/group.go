package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/v3lichko/student-distribution/internal/models"
	"github.com/v3lichko/student-distribution/internal/response"
)

type GroupHandler struct {
	db *pg.DB
}

func NewGroupHandler(db *pg.DB) *GroupHandler {
	return &GroupHandler{
		db: db,
	}
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

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	json.NewDecoder(r.Body).Decode(&group)
	h.db.Model(&group).Insert()
	response.WriteJSON(w, http.StatusCreated, group)
}

func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	groups := make([]models.Group, 0)
	h.db.Model(&groups).Select()
	response.WriteJSON(w, http.StatusOK, groups)
}
