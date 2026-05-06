package exceptions

const (
	ErrProductNotFound      string = "errProductNotFound"
	ErrProductNameRequired  string = "errProductNameRequired"
	ErrProductPriceInvalid  string = "errProductPriceInvalid"
	ErrProductStatusInvalid string = "errProductStatusInvalid"

	ErrOnInsertProduct       string = "errOnInsertProduct"
	ErrOnUpdateProduct       string = "errOnUpdateProduct"
	ErrOnDeleteProduct       string = "errOnDeleteProduct"
	ErrOnFindProductByID     string = "errOnFindProductByID"
	ErrOnPaginateProducts    string = "errOnPaginateProducts"
	ErrOnPublishProductEvent string = "errOnPublishProductEvent"
)
