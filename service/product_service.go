package service

import (
	"errors"
	"prod_app/domain"
	"prod_app/persistance"
	"prod_app/persistance/constants/errmsgs"
	"prod_app/service/model"
)

type IProductService interface {
	Add(productCreate model.ProductCreate) error
	DeleteById(productId int64) error
	GetById(productId int64) (domain.Product, error)
	UpdatePrice(productId int64, newPrice float32) error
	GetAllProducts() []domain.Product
	GetAllProductsByStore(storeName string) []domain.Product
}

type ProductService struct {
	productRepository persistance.IProductRepository
}

func NewProductService(productRepository persistance.IProductRepository) IProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (productService *ProductService) Add(productCreate model.ProductCreate) error {
	validationErr := validateProductCreate(productCreate)
	if validationErr != nil {
		return validationErr
	}
	return productService.productRepository.AddProduct(
		domain.Product{
			Name:     productCreate.Name,
			Price:    productCreate.Price,
			Discount: productCreate.Discount,
			Store:    productCreate.Store,
		},
	)
}

func (productService *ProductService) DeleteById(productId int64) error {
	if productId <= 0 {
		return errors.New(
			errmsgs.ServiceIdInvalid,
		)
	}
	return productService.productRepository.DeleteById(productId)
}

func (productService *ProductService) GetById(productId int64) (domain.Product, error) {
	if productId <= 0 {
		return domain.Product{}, errors.New(
			errmsgs.ServiceIdInvalid,
		)
	}
	return productService.productRepository.GetById(productId)
}

func (productService *ProductService) UpdatePrice(productId int64, newPrice float32) error {
	if productId <= 0 {
		return errors.New(
			errmsgs.ServiceIdInvalid,
		)
	}
	if newPrice <= 0 {
		return errors.New(
			errmsgs.ServiceNewPriceInvalid,
		)
	}
	return productService.productRepository.UpdatePrice(productId, newPrice)
}

func (productService *ProductService) GetAllProducts() []domain.Product {
	return productService.productRepository.GetAllProducts()
}

func (productService *ProductService) GetAllProductsByStore(storeName string) []domain.Product {
	return productService.productRepository.GetAllProductsByStore(storeName)
}

func validateProductCreate(productCreate model.ProductCreate) error {
	if productCreate.Discount > 70.0 {
		return errors.New("Discount can not be bigger than 70.0")
	}
	return nil
}
