package main

import (
	"context"
	"prod_app/common/app"
	"prod_app/common/postgresql"
	"prod_app/controller"
	"prod_app/persistance"
	"prod_app/service"

	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()
	configuration := app.NewConfigurationManager()
	dbpool := postgresql.GetConnectionPool(ctx, configuration.PostgreSqlConnectionString)
	productRepo := persistance.NewProductRepository(dbpool)
	productService := service.NewProductService(productRepo)
	productController := controller.NewnProductController(productService)

	ec := echo.New()
	productController.RegisterRoutes(ec)

	ec.Start("localhost:8080")
}
