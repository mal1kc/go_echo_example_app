package controller

import (
	"net/http"
	"prod_app/controller/request"
	"prod_app/controller/response"
	"prod_app/service"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ProductController struct {
	productService service.IProductService
}

func NewnProductController(productService service.IProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (productController *ProductController) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/v1/products", productController.GetAllProducts)
	e.POST("/api/v1/products", productController.AddProduct)
	e.GET("/api/v1/products/:id", productController.GetProductById)
	e.PUT("/api/v1/products/:id", productController.UpdatePrice)
	e.DELETE("/api/v1/products/:id", productController.DeleteProductById)
}

func (productController *ProductController) GetAllProducts(c echo.Context) error {
	store := c.QueryParam("store")
	if len(store) == 0 {
		products := productController.productService.GetAllProducts()
		return c.JSON(http.StatusOK, response.ToResponses(products))
	}
	productsWithStoreFilter := productController.productService.GetAllProductsByStore(store)
	return c.JSON(http.StatusOK, response.ToResponses(productsWithStoreFilter))
}

func (productController *ProductController) AddProduct(c echo.Context) error {
	productRequest := new(request.AddProductRequest)
	bindErr := c.Bind(productRequest)
	if bindErr != nil {
		log.Errorf("this shit came %v\n", bindErr.Error())
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{ErrorDescription: "Bad data as json check keys etc"})
	}
	addErr := productController.productService.Add(productRequest.ToModel())
	if addErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorDescription: "Given data is not valid",
		})
	}
	return c.NoContent(http.StatusCreated)
}

func (productController *ProductController) GetProductById(c echo.Context) error {
	strId := c.Param("id")
	prodId, validatErr := strconv.Atoi(strId)
	if validatErr != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: "cant't parse id",
		})
	}

	product, servErr := productController.productService.GetById(int64(prodId))
	if servErr != nil {
		return c.JSON(
			http.StatusNotFound, response.ErrorResponse{
				ErrorDescription: servErr.Error(),
			},
		)
	}
	return c.JSON(http.StatusOK, response.ToResponse(product))
}

func (productController *ProductController) UpdatePrice(c echo.Context) error {
	strId := c.Param("id")
	prodId, conv_err := strconv.Atoi(strId)
	if conv_err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{ErrorDescription: "cant't parse id"})
	}

	newPriceStr := c.QueryParam("newPrice")
	if len(newPriceStr) == 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: "newPrice param is required.",
		})
	}
	newPrice, parseErr := strconv.ParseFloat(newPriceStr, 32)
	if parseErr != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{ErrorDescription: "can't parse newPrice param"})
	}

	servErr := productController.productService.UpdatePrice(int64(prodId), float32(newPrice))
	if servErr != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: servErr.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}
func (productController *ProductController) DeleteProductById(c echo.Context) error {
	strId := c.Param("id")
	prodId, conv_err := strconv.Atoi(strId)
	if conv_err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{ErrorDescription: "cant't parse id"})
	}

	servErr := productController.productService.DeleteById(int64(prodId))
	if servErr != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			ErrorDescription: servErr.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}
