package infrastructure

import (
	"context"
	"fmt"
	"os"
	"prod_app/common/app"
	"prod_app/common/postgresql"
	"prod_app/domain"
	"prod_app/persistance"
	"prod_app/persistance/constants/errmsgs"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

var productRepository persistance.IProductRepository
var dbPool *pgxpool.Pool
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()
	test_configman := app.NewConfigurationManager()
	dbPool = postgresql.GetConnectionPool(ctx, test_configman.PostgreSqlConnectionString)

	productRepository = persistance.NewProductRepository(dbPool)

	fmt.Println("Before all tests")
	clear(ctx, dbPool)
	exitCode := m.Run()
	fmt.Println("After all tests")
	os.Exit(exitCode)
}

func setup(ctx context.Context, dbPool *pgxpool.Pool) {
	TestDataInitialize(ctx, dbPool)
}
func clear(ctx context.Context, dbPool *pgxpool.Pool) {
	TruncateTestData(ctx, dbPool)
}

func TestGetAllProducts(t *testing.T) {
	t.Run("GetAllProductsEmpty", func(t *testing.T) {
		actualProducts := productRepository.GetAllProducts()
		assert.Equal(t, 0, len(actualProducts))
	})

	setup(ctx, dbPool)

	expectedProducts := []domain.Product{
		{
			Id:       1,
			Name:     "AirFryer",
			Price:    3000.0,
			Discount: 22.0,
			Store:    "ABC TECH",
		},
		{
			Id:       2,
			Name:     "Ütü",
			Price:    1500.0,
			Discount: 10.0,
			Store:    "ABC TECH",
		},
		{
			Id:       3,
			Name:     "Çamaşır Makinesi",
			Price:    10000.0,
			Discount: 15.0,
			Store:    "ABC TECH",
		},
		{
			Id:       4,
			Name:     "Lambader",
			Price:    2000.0,
			Discount: 0.0,
			Store:    "Dekorasyon Sarayı",
		},
	}
	t.Run("GetAllProductsFull", func(t *testing.T) {
		actualProducts := productRepository.GetAllProducts()
		assert.Equal(t, 4, len(actualProducts))
		assert.Equal(t, expectedProducts, actualProducts)
		assert.NotEqual(t, expectedProducts[0], actualProducts[1])
	})

	clear(ctx, dbPool)
}

func TestGetAllProductsByStore(t *testing.T) {
	t.Run("TestGetAllProductsByStoreEmpty", func(t *testing.T) {
		actualProducts := productRepository.GetAllProductsByStore("")
		assert.Equal(t, 0, len(actualProducts))
	})

	setup(ctx, dbPool)

	expectedProducts := []domain.Product{
		{
			Id:       1,
			Name:     "AirFryer",
			Price:    3000.0,
			Discount: 22.0,
			Store:    "ABC TECH",
		},
		{
			Id:       2,
			Name:     "Ütü",
			Price:    1500.0,
			Discount: 10.0,
			Store:    "ABC TECH",
		},
		{
			Id:       3,
			Name:     "Çamaşır Makinesi",
			Price:    10000.0,
			Discount: 15.0,
			Store:    "ABC TECH",
		},
	}
	t.Run("TestGetAllProductsByStoreFull", func(t *testing.T) {
		actualProducts := productRepository.GetAllProductsByStore("ABC TECH")
		assert.Equal(t, 3, len(actualProducts))
		assert.Equal(t, expectedProducts, actualProducts)
		assert.NotEqual(t, expectedProducts[0], actualProducts[1])
	})

	clear(ctx, dbPool)
}

func TestAddProduct(t *testing.T) {
	setup(ctx, dbPool)

	expectedProducts := []domain.Product{
		{
			Id:       5,
			Name:     "some thing in the way",
			Price:    5555.0,
			Discount: 69.0,
			Store:    "The batman Store",
		},
	}
	t.Run("TestAddProduct", func(t *testing.T) {
		assert.Equal(t, nil,
			productRepository.AddProduct(expectedProducts[0]),
		)
		actualProducts := productRepository.GetAllProductsByStore("The batman Store")
		assert.Equal(t, 1, len(actualProducts))
		assert.Equal(t, expectedProducts, actualProducts)
	})

	clear(ctx, dbPool)
}

func TestGetProductById(t *testing.T) {
	setup(ctx, dbPool)
	expectedProduct := domain.Product{
		Id:       1,
		Name:     "AirFryer",
		Price:    3000.0,
		Discount: 22.0,
		Store:    "ABC TECH",
	}
	t.Run("GetProductById", func(t *testing.T) {
		actualProduct, _ := productRepository.GetById(1)
		_, acterr := productRepository.GetById(5)

		assert.Equal(t, expectedProduct, actualProduct)
		assert.Equal(t, fmt.Sprintf(errmsgs.BY_ID_NOT_FOUND_FMT, 5), acterr.Error())
	})
	clear(ctx, dbPool)
}

func TestDeleteProductById(t *testing.T) {
	setup(ctx, dbPool)
	t.Run("DeleteById", func(t *testing.T) {
		actualErr := productRepository.DeleteById(2)
		assert.Equal(t, nil, actualErr)

		actualErr = productRepository.DeleteById(9)
		assert.Equal(t, fmt.Sprintf(errmsgs.BY_ID_NOT_FOUND_FMT, 9), actualErr.Error())
	})

	t.Run("afterDeleteByIdGetAllProducts", func(t *testing.T) {
		actualProducts := productRepository.GetAllProducts()
		assert.Equal(t, 3, len(actualProducts))
	})

	clear(ctx, dbPool)
}

func TestUpdateProductPrice(t *testing.T) {
	setup(ctx, dbPool)
	t.Run("UpdateProductPrice", func(t *testing.T) {
		uptdPrcTarget := float32(2000.0)
		updtIdTarget := int64(1)
		expectedProduct, _ := productRepository.GetById(updtIdTarget)
		assert.NotEqual(t, expectedProduct.Price, uptdPrcTarget)
		actualErr := productRepository.UpdatePrice(updtIdTarget, uptdPrcTarget)
		assert.Equal(t, nil, actualErr)

		// check is it updated correctly
		expectedProduct, _ = productRepository.GetById(updtIdTarget)
		assert.Equal(t, expectedProduct.Price, uptdPrcTarget)

		// non -existed id
		actualErr = productRepository.UpdatePrice(9, uptdPrcTarget)
		assert.Equal(t, fmt.Sprintf(errmsgs.BY_ID_NOT_FOUND_FMT, 9), actualErr.Error())
	})

	clear(ctx, dbPool)
}
