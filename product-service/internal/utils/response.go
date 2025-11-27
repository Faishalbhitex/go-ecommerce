package utils

import (
	"encoding/json"
	"net/http"
	"product-service/internal/dto"
	"product-service/internal/models"
)

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func Err(w http.ResponseWriter, appErr *AppError) {
	JSON(w, appErr.Status, map[string]any{
		"error": map[string]string{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	})
}

func ToProductResponse(p *models.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Qty:         p.Qty,
		Category:    p.Category,
		CreateAt:    p.CreateAt,
		UpdateAt:    p.UpdateAt,
	}
}

func ToProductResponseSlice(products []*models.Product) []dto.ProductResponse {
	list := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		list[i] = ToProductResponse(p)
	}
	return list
}
