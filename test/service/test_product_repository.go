package service

import (
	"fmt"
	"prod_app/domain"
	"prod_app/persistance"
	"prod_app/persistance/constants/errmsgs"
)

type TestProductRepository struct {
	products []domain.Product
}

func NewTestProductRepository(initialProducts []domain.Product) persistance.IProductRepository {
	return &TestProductRepository{
		products: initialProducts,
	}
}

func (productRepository *TestProductRepository) GetAllProducts() []domain.Product {
	return productRepository.products
}
func (productRepository *TestProductRepository) GetAllProductsByStore(storeName string) []domain.Product {
	res_products := []domain.Product{}
	for _, pr := range productRepository.products {
		if pr.Store == storeName {
			res_products = append(res_products, pr)
		}
	}
	return res_products
}
func (productRepository *TestProductRepository) AddProduct(product domain.Product) error {
	product.Id = int64(len(productRepository.products)) + 1
	productRepository.products = append(productRepository.products, product)
	return nil
}
func (productRepository *TestProductRepository) GetById(productId int64) (domain.Product, error) {
	index := findProdIndexById(productRepository.products, productId)
	if index == -1 {
		return domain.Product{}, fmt.Errorf(errmsgs.BY_ID_NOT_FOUND_FMT, productId)
	}
	pr := productRepository.products[index]
	return pr, nil
}

func (productRepository *TestProductRepository) DeleteById(productId int64) error {
	index := findProdIndexById(productRepository.products, productId)
	if index == -1 {
		return fmt.Errorf(errmsgs.BY_ID_NOT_FOUND_FMT, productId)
	}
	productRepository.products = removeProdFromSlice(productRepository.products, index)
	return nil
}

func (productRepository *TestProductRepository) UpdatePrice(productId int64, newPrice float32) error {
	index := findProdIndexById(productRepository.products, productId)
	if index == -1 {
		return fmt.Errorf(errmsgs.BY_ID_NOT_FOUND_FMT, productId)
	}
	pr := productRepository.products[index]
	pr.Price = newPrice
	productRepository.products[index] = pr
	return nil
}

func findProdIndexById(productSlice []domain.Product, productId int64) int {
	var index int = -1
	if int64(len(productSlice)) > productId {
		for i, pr := range productSlice {
			if pr.Id == productId {
				index = i
				break
			}
		}
	}
	return index
}

func removeProdFromSlice(s []domain.Product, i int) []domain.Product {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
