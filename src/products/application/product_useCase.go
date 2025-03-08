package application

import "consumer_api/src/products/domain"


type ProductProcessingUseCase struct {
	consumer domain.ProductConsumer
}

func NewProductProcessingUseCase(consumer domain.ProductConsumer) *ProductProcessingUseCase {
	return &ProductProcessingUseCase{consumer: consumer}
}

func (uc *ProductProcessingUseCase) ProcessProduct(product *domain.Product) error {
	// Aquí podemos hacer el procesamiento, como actualizar la base de datos, realizar validaciones, etc.
	return uc.consumer.ProcessProduct(product)
}