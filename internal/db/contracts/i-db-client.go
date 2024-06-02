package contracts

type IDbClient interface {
	BulkInsertFromCsv(query string) error
	CloseConnection()
}
