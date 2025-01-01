package errmsgs

type error_fmt_msg string

const (
	BY_ID_ERR_FMT       error_fmt_msg = "Error while getting product with id %d"
	BY_ID_NOT_FOUND_FMT               = "Product with id %d not found"
	DEL_BY_ID_ERR_FMT                 = "Error while deleting product with id %d"
	UPDT_PRICE_ERR_FMT                = "Error while updating product price with id : %d"
)
