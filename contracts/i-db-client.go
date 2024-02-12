package contracts

type IDbClient interface {
	InsertRows(query string) error
	CloseConnection()
}
