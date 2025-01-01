package service

import (
	"os"
	"prod_app/domain"
	"prod_app/service"
	"prod_app/service/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

var prductService service.IProductService
var initial_db_data []domain.Product

func TestMain(m *testing.M) {
	initial_db_data = []domain.Product{
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
		}}

	prductService = service.NewProductService(NewTestProductRepository(initial_db_data))
	eCode := m.Run()
	os.Exit(eCode)
}

func Test_ShouldGetAllProducts(t *testing.T) {
	t.Run(
		"ShouldGetAllProducts", func(t *testing.T) {
			expected_data := initial_db_data
			actual_data := prductService.GetAllProducts()
			assert.Equal(t, expected_data, actual_data)
		},
	)
}

func Test_ShouldGetAllProductsByStore(t *testing.T) {
	t.Run(
		"ShouldGetAllProductsByStore", func(t *testing.T) {
			expected_data := initial_db_data[:len(initial_db_data)-1]
			actual_data := prductService.GetAllProductsByStore("ABC TECH")
			assert.NotEqual(t, initial_db_data, actual_data)
			assert.Equal(t, expected_data, actual_data)
		},
	)
}

func Test_WhenThereIsNoStoreProduct_ShouldGetEmptyProductSliceByStore(t *testing.T) {
	t.Run(
		"WhenThereIsNoStoreProduct_ShouldGetEmptyProductSliceByStore", func(t *testing.T) {
			actual_data := prductService.GetAllProductsByStore("aaaaa")
			assert.NotEqual(t, initial_db_data, actual_data)
			assert.Equal(t, []domain.Product{}, actual_data)
		},
	)
}
func Test_WhenDiscountIsHigherThan70_ShouldNotAddProduct(t *testing.T) {
	t.Run("WhenDiscountIsHigherThan70_ShouldNotAddProduct", func(t *testing.T) {
		expected_err := `Discount can not be bigger than 70.0`
		created_prod := model.ProductCreate{
			Name:     "Aa",
			Price:    200.0,
			Discount: 90.0,
			Store:    "aa store",
		}
		actual_err := prductService.Add(created_prod)
		assert.Equal(t, expected_err, actual_err.Error())
	})
}

func Test_WhenValidationErrorOccured_ShouldNotAddProduct(t *testing.T) {
	t.Run("WhenValidationErrorOccured_ShouldNotAddProduct", func(t *testing.T) {
		expected_err := `Discount can not be bigger than 70.0`
		created_prod := model.ProductCreate{
			Name:     "Aa",
			Price:    200.0,
			Discount: 90.0,
			Store:    "aa store",
		}
		actual_err := prductService.Add(created_prod)
		assert.Equal(t, expected_err, actual_err.Error())
	})
}

func Test_WhenNoValidationErrorOccured_ShouldAddProduct(t *testing.T) {
	t.Run("WhenNoValidationErrorOccured_ShouldAddProduct", func(t *testing.T) {
		created_prod := model.ProductCreate{
			Name:     "Aa",
			Price:    200.0,
			Discount: 30.0,
			Store:    "aa store",
		}
		actual_res := prductService.Add(created_prod)
		assert.Equal(t, nil, actual_res)
		expected_product := domain.Product{
			Id:       int64(len(initial_db_data)) + 1,
			Name:     created_prod.Name,
			Price:    created_prod.Price,
			Discount: created_prod.Discount,
			Store:    created_prod.Store,
		}
		pructs := prductService.GetAllProducts()
		actual_product := pructs[len(pructs)-1]

		assert.Equal(t, expected_product, actual_product)
	})
}
