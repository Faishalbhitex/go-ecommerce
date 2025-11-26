package dto

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	Qty         int     `json:"qty"`
	Category    string  `json:"category,omitempty"`
}

type UpdateProductRequest = CreateProductRequest

type PatchProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Qty         *int     `json:"qty,omitempty"`
	Category    *string  `json:"category,omitempty"`
}
