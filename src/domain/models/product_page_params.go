package models

type ProductPageParams struct {
	Page     int    `query:"page" validate:"required,min=1"`
	PageSize int    `query:"pageSize" validate:"required,min=1,max=100"`
	Status   string `query:"status" validate:"omitempty,oneOfProductStatus"`
}

type ProductPage struct {
	Items      []Product `json:"items"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalItems int64     `json:"total_items"`
}
