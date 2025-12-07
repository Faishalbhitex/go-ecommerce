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
		utils.Err(w, utils.BadRequest("invalid id"))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	p, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.Err(w, utils.NotFound("product not found"))
			return
		}
		h.errLog.Println("getbyid:", err)
		utils.Err(w, utils.Internal("server error"))
		return
	}
	utils.JSON(w, http.StatusOK, utils.ToProductResponse(p))
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	list, err := h.svc.List(ctx)
	if err != nil {
		h.errLog.Println("List:", err)
		utils.Err(w, utils.Internal("failed to list products"))
		return
	}
	utils.JSON(w, http.StatusOK, map[string]any{
		"data":  utils.ToProductResponseSlice(list),
		"count": len(list),
	})
}

func (h *ProductHandler) ListPaged(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p >= 1 {
		page = p
	}
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l >= 1 && l <= 100 {
		limit = l
	}
	products, err := h.svc.ListPaged(r.Context(), page, limit)
	if err != nil {
		h.errLog.Println("list paged:", err)
		utils.Err(w, utils.Internal("failed to get paged products"))
		return
	}

	resp := map[string]any{
		"data":  utils.ToProductResponseSlice(products),
		"page":  page,
		"limit": limit,
	}
	utils.JSON(w, http.StatusOK, resp)
}

func (h *ProductHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		utils.Err(w, utils.BadRequest("query parameter 'q' is required"))
		return
	}

	products, err := h.svc.Search(r.Context(), q)
	if err != nil {
		h.errLog.Println("search:", err)
		utils.Err(w, utils.Internal("search failed"))
		return
	}
	utils.JSON(w, http.StatusOK, map[string]any{
		"data": utils.ToProductResponseSlice(products),
	})
}

//// @HANDLER:READ-END

// // @HANDLER:UPDATE-BEGIN
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.errLog.Println("update:", err)
		utils.Err(w, utils.BadRequest("invalid id"))
		return
	}
	var input dto.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Err(w, utils.BadRequest("invalid JSON payload"))
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
		utils.Err(w, utils.Validation("qty cannot be negative"))
		return
	}

	p := &models.Product{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Qty:         input.Qty,
		Category:    input.Category,
	}

	updated, err := h.svc.Update(ctx, p)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.Err(w, utils.NotFound("product not found"))
			return
		}
		h.errLog.Println("update product:", err)
		utils.Err(w, utils.Internal("failed to update product"))
		return
	}

	utils.JSON(w, http.StatusOK, utils.ToProductResponse(updated))
}

func (h *ProductHandler) Patch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.errLog.Println("patch:", err)
		utils.Err(w, utils.BadRequest("invalid id"))
		return
	}

	var input dto.PatchProductRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Err(w, utils.BadRequest("invalid JSON payload"))
		return
	}

	if input.Price != nil && *input.Price <= 0 {
		utils.Err(w, utils.Validation("price must be greater than 0"))
		return
	}
	if input.Qty != nil && *input.Qty < 0 {
		utils.Err(w, utils.Validation("qty cannot be negative"))
		return
	}

	patch := &models.ProductPatch{}
	if input.Name != nil {
		patch.Name = input.Name
	}
	if input.Description != nil {
		patch.Description = input.Description
	}
	if input.Price != nil {
		patch.Price = input.Price
	}
	if input.Qty != nil {
		patch.Qty = input.Qty
	}
	if input.Category != nil {
		patch.Category = input.Category
	}

	updated, err := h.svc.Patch(ctx, id, patch)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.Err(w, utils.NotFound("product not found"))
			return
		}
		h.errLog.Println("patch product:", err)
		utils.Err(w, utils.Internal("failed to patch product"))
		return
	}

	utils.JSON(w, http.StatusOK, utils.ToProductResponse(updated))
}

//// @HANDLER:UPDATE-END

// // @HANDLER:DELETE-BEGIN
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Err(w, utils.BadRequest("invalid id"))
		return
	}
	if err := h.svc.Delete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			utils.Err(w, utils.NotFound("product not found"))
			return
		}
		h.errLog.Println("delete:", err)
		utils.Err(w, utils.Internal("failed to delete product"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//// @HANDLER:DELETE-END
