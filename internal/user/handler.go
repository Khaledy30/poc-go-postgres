package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"uuid"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "body inválido.")
		return
	}

	user, err := h.svc.Create(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respond(w, http.StatusCreated, user)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	users, err := h.svc.List(r.Context(), page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao listar usuários")
		return
	}

	if users == nil {
		users = []User{}
	}

	respond(w, http.StatusOK, users)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if _, err := uuid.Parse(idStr); err != nil {
		respondError(w, http.StatusBadRequest, "id inválido")
		return
	}

	user, err := h.svc.GetByID(r.Context(), idStr)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao buscar usuário")
		return
	}
	if user == nil {
		respondError(w, http.StatusNotFound, "usuário não encontrado")
		return
	}

	respond(w, http.StatusOK, user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if _, err := uuid.Parse(idStr); err != nil {
		respondError(w, http.StatusBadRequest, "id inválido")
		return
	}

	err := h.svc.Delete(r.Context(), idStr)
	switch {
	case errors.Is(err, ErrNotFound):
		respondError(w, http.StatusNotFound, "usuário não encontrado")
		return
	case errors.Is(err, ErrAlreadyDeleted):
		respond(w, http.StatusOK, map[string]string{
			"message": "usuário já foi deletado",
		})
		return
	case err != nil:
		respondError(w, http.StatusInternalServerError, "erro ao deletar usuário")
		return
	}

	respond(w, http.StatusOK, map[string]string{
		"message": "usuário deletado com sucesso",
	})
}

func respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respond(w, status, map[string]string{"error": message})
}
