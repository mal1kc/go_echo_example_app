package persistance

import (
	"context"
	"fmt"
	"prod_app/domain"
	"prod_app/persistance/constants/errmsgs"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

type IProductRepository interface {
	GetAllProducts() []domain.Product
	GetAllProductsByStore(storeName string) []domain.Product
	AddProduct(product domain.Product) error
	GetById(productId int64) (domain.Product, error)
	DeleteById(productId int64) error
	UpdatePrice(productId int64, newPrice float32) error
}

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) IProductRepository {
	return &ProductRepository{
		dbPool: dbPool,
	}
}

func extractProductsFromRows(rows pgx.Rows) []domain.Product {
	products := []domain.Product{}
	tmpprod := domain.Product{}

	for rows.Next() {
		rErr := rows.Scan(&tmpprod.Id, &tmpprod.Name, &tmpprod.Price, &tmpprod.Discount, &tmpprod.Store)
		if rErr != nil {
			log.Error(string(errmsgs.LOG_EXTRACT_FROM_ROWS_ERR), rows.Err())
		}
		products = append(products, tmpprod)
		// this not be needed because append operation cretes new copy because products is not pointer slice
		// tmpprod = domain.Product{}
	}
	return products
}

func (productRepository *ProductRepository) GetAllProducts() []domain.Product {
	ctx := context.Background()
	rows, qErr := productRepository.dbPool.Query(ctx, "select * from products")
	if qErr != nil {
		log.Error(string(errmsgs.LOG_GETALL_ERR), qErr)
		return []domain.Product{}
	}

	return extractProductsFromRows(rows)
}

func (productRepository *ProductRepository) GetAllProductsByStore(storeName string) []domain.Product {
	ctx := context.Background()
	getProductsbyStoreNameSQL := `select * from products where store = $1`
	productRows, qErr := productRepository.dbPool.Query(ctx, getProductsbyStoreNameSQL, storeName)
	if qErr != nil {
		log.Error(string(errmsgs.LOG_GETALL_BYSTORE_ERR), qErr)
		return []domain.Product{}
	}

	return extractProductsFromRows(productRows)
}

func (productRepository *ProductRepository) AddProduct(product domain.Product) error {
	ctx := context.Background()
	insert_sql := `Insert into products (name,price,discount,store) VALUES( $1, $2, $3, $4)`
	cmTag, err := productRepository.dbPool.Exec(ctx, insert_sql,
		product.Name,
		product.Price,
		product.Discount,
		product.Store,
	)
	if err != nil {
		log.Error(string(errmsgs.LOG_ADD_PRODUCT_ERR), err)
		return err
	}

	log.Info(fmt.Sprintf("Product added with %v", cmTag))
	return nil
}

func (productRepository *ProductRepository) GetById(productId int64) (domain.Product, error) {
	ctx := context.Background()
	getbyid_sql := `select * from products where id = $1`
	qRow := productRepository.dbPool.QueryRow(ctx, getbyid_sql, productId)
	prod := domain.Product{}

	scErr := qRow.Scan(&prod.Id, &prod.Name, &prod.Price, &prod.Discount, &prod.Store)
	if scErr != nil {
		if scErr.Error() == string(errmsgs.NOT_FOUND) {
			return prod, fmt.Errorf(string(errmsgs.BY_ID_NOT_FOUND_FMT), productId)
		}
		return prod, fmt.Errorf(string(errmsgs.BY_ID_ERR_FMT), productId)
	}
	return prod, nil
}

func (productRepository *ProductRepository) DeleteById(productId int64) error {
	ctx := context.Background()
	_, qErr := productRepository.GetById(productId)
	if qErr != nil {
		return qErr
	}

	deleteSql := `delete from products where id = $1`
	_, execErr := productRepository.dbPool.Exec(ctx, deleteSql, productId)
	if execErr != nil {
		log.Error(errmsgs.LOG_DEL_PRODUCT_BYID_ERR, execErr.Error())
		return fmt.Errorf(string(errmsgs.DEL_BY_ID_ERR_FMT), productId)
	}
	log.Info("prdocut deleted id: ", productId)
	return nil
}

func (productRepository *ProductRepository) UpdatePrice(productId int64, newPrice float32) error {
	ctx := context.Background()
	updateSql := `Update products set price = $1 where id = $2`
	res, execErr := productRepository.dbPool.Exec(ctx, updateSql, newPrice, productId)
	if execErr != nil {
		log.Error(errmsgs.LOG_UPDT_PRICE_ERR, execErr.Error())
		return fmt.Errorf(string(errmsgs.UPDT_PRICE_ERR_FMT), productId)
	}
	if res.RowsAffected() == 0 {
		log.Info(errmsgs.LOG_UPDT_PRICE_NOEFFECT_ERR)
		return fmt.Errorf(string(errmsgs.BY_ID_NOT_FOUND_FMT), productId)
	}
	return nil
}
