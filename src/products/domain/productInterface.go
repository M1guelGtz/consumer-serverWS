package domain

type ProductConsumer interface {
	ProcessProduct(product *Product) error
}
type Consumer interface {
    Consume() error
}