package pack

import (
	"encoding/json"
	"net/http"
	"reparttask/storage"
	"reparttask/utils"
	"strconv"
)

type SizePayload struct {
	Sizes []int `json:"sizes"`
}

type Handler struct {
	db storage.Storage
}

func NewHandler(db storage.Storage) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /pack", h.handleAddPacks)
	router.HandleFunc("DELETE /pack/{size}", h.handleRemovePack)
	router.HandleFunc("DELETE /packs", h.handleRemovePacks)
}

func (h *Handler) handleAddPacks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var pk SizePayload
	err := json.NewDecoder(r.Body).Decode(&pk)
	if err != nil {
		utils.WriteOutput(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if len(pk.Sizes) == 0 {
		utils.WriteOutput(w, http.StatusBadRequest, map[string]string{"error": "pack size must be positive"})
		return
	}

	err = h.db.AddPacks(pk.Sizes)
	if err != nil {
		utils.WriteOutput(w, http.StatusInternalServerError, map[string]string{"error": "an error has occurred"})
		return
	}

	utils.WriteOutput(w, http.StatusCreated, map[string]string{"status": "success"})
}

func (h *Handler) handleRemovePack(w http.ResponseWriter, r *http.Request) {
	size := r.PathValue("size")
	if size == "" {
		utils.WriteOutput(w, http.StatusBadRequest, map[string]string{"error": "you must provide a size value"})
		return
	}

	nr, err := strconv.Atoi(size)
	if err != nil {
		utils.WriteOutput(w, http.StatusBadRequest, map[string]string{"error": "please provide a numeric value"})
		return
	}

	if nr <= 0 {
		utils.WriteOutput(w, http.StatusBadRequest, map[string]string{"error": "you must provide a positive value"})
		return
	}

	err = h.db.RemovePack(nr)
	if err != nil {
		utils.WriteOutput(w, http.StatusInternalServerError, map[string]string{"error": "an error has occurred"})
		return
	}

	utils.WriteOutput(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *Handler) handleRemovePacks(w http.ResponseWriter, r *http.Request) {
	h.db.RemovePacks()
	utils.WriteOutput(w, http.StatusOK, map[string]string{"status": "success"})
}
