package response

import "prod_app/domain"

type ErrorResponse struct {
	ErrorDescription string `json:"error_description"`
}

type ProductResponse struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Discount float32 `json:"discount"`
	Store    string  `json:"store"`
}

func ToResponse(product domain.Product) ProductResponse {
	return ProductResponse{
		Name:     product.Name,
		Price:    product.Price,
		Discount: product.Discount,
		Store:    product.Store,
	}
}

func ToResponses(productSlice []domain.Product) []ProductResponse {
	response_slice := []ProductResponse{}
	for _, product := range productSlice {
		response_slice = append(response_slice, ProductResponse{
			Name:     product.Name,
			Price:    product.Price,
			Discount: product.Discount,
			Store:    product.Store,
		})
	}
	return response_slice
}
