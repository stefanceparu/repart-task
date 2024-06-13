package order

import (
	"net/http"
	"reparttask/service"
	"reparttask/storage"
	"reparttask/utils"
	"strconv"
)

type Handler struct {
	db   storage.Storage
	calc service.Calculator
}

func NewHandler(db storage.Storage, calc service.Calculator) *Handler {
	return &Handler{db: db, calc: calc}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /order/{items}", h.handleGetOrder)
}

func (h *Handler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	items := r.PathValue("items")
	if items == "" {
		utils.WriteOutput(w, http.StatusBadRequest, map[string]string{"error": "you must provide a number of items"})
		return
	}

	nr, err := strconv.Atoi(items)
	if err != nil {
		utils.WriteOutput(w, http.StatusBadRequest, map[string]string{"error": "please provide a numeric value"})
		return
	}

	if nr <= 0 {
		utils.WriteOutput(w, http.StatusBadRequest, map[string]string{"error": "please provide a number greater than zero"})
		return
	}

	packs := h.db.GetPacks()
	if len(packs) == 0 {
		utils.WriteOutput(w, http.StatusBadRequest, map[string]string{"error": "you must first add some packaging sizes"})
		return
	}

	result := h.calc.CalculatePacks(packs, nr)

	utils.WriteOutput(w, http.StatusOK, result)
}
