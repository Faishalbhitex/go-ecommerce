package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"product-service/internal/dto"
	"product-service/internal/models"
	"product-service/internal/repository"
	"product-service/internal/service"
	"product-service/internal/utils"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	svc     service.ProductService
	infoLog *log.Logger
	errLog  *log.Logger
}

// Constructorr
func NewProductHandler(svc service.ProductService, info, errl *log.Logger) *ProductHandler {
	return &ProductHandler{
		svc:     svc,
		infoLog: info,
		errLog:  errl,
	}
}

// // @HANDLER:CREATE-BEGIN
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Err(w, utils.BadRequest("Invalid JSON payload"))
		return
	}

	if input.Name == "" {
		utils.Err(w, utils.Validation("name is required"))
		return
	}

	if input.Price <= 0 {
		utils.Err(w, utils.Validation("price must be greater than 0"))
		return
	}
	if input.Qty < 0 {
		utils.Err(w, utils.Validation("qty cannot br negative"))
		return
	}

	p := &models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Qty:         input.Qty,
		Category:    input.Category,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.svc.Create(ctx, p); err != nil {
		h.errLog.Println("create product:", err)
		utils.Err(w, utils.Internal("Failed to create product"))
		return
	}

	resp := dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Qty:         p.Qty,
		Category:    p.Category,
		CreateAt:    p.CreateAt,
		UpdateAt:    p.UpdateAt,
	}

	utils.JSON(w, http.StatusCreated, resp)
}

//// @HANDLER:CREATE-END

// // @HANDLER:READ-BEGIN
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	p, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.errLog.Println("getbyid:", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	list, err := h.svc.List(ctx)
	if err != nil {
		h.errLog.Println("List:", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *ProductHandler) ListPaged(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	p, err := h.svc.ListPaged(r.Context(), page, limit)
	if err != nil {
		h.errLog.Println("list paged:", err)
		http.Error(w, "failed to list paged products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	p, err := h.svc.Search(r.Context(), q)
	if err != nil {
		h.errLog.Println("search:", err)
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

//// @HANDLER:READ-END

// // @HANDLER:UPDATE-BEGIN
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.errLog.Println("update:", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	p.ID = id
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.Update(ctx, &p); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) Patch(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.errLog.Println("patch:", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var patch models.ProductPatch
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.svc.Patch(r.Context(), id, &patch); err != nil {
		http.Error(w, "patch failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "patched"}`))
}

//// @HANDLER:UPDATE-END

// // @HANDLER:DELETE-BEGIN
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.Delete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.errLog.Println("delete:", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//// @HANDLER:DELETE-END
