package contracts

type IMessageConsumer interface {
	ConsumeMessage() ([]byte, error)
	CloseConnection() error
}
