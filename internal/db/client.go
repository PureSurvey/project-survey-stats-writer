package db

import "io"

type Client interface {
	BulkInsert(inputStream io.Reader, tableName string) error
	CloseConnection()
}
