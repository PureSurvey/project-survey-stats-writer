package db

type Client interface {
	BulkInsertFromCsv(query string) error
	CloseConnection()
}
