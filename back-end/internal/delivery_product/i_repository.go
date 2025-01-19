package delivery_product

type IDeliveryProductRepository interface {
	// Returns the created deliveryProduct
	Create(DeliveryProductParams) (DeliveryProductModel, error)

	// Returns the found deliveryProducts
	// NOTE: reusing the same type for a filter and a "constructor" is not
	// ideal at all, but it will save on code repetition
	GetAll(filter DeliveryProductParams) ([]DeliveryProductModel, error)

	// Returns the found deliveryProduct
	GetOne(id string) (DeliveryProductModel, error)

	// Returns amount of deleted deliveryProduct
	DeleteOne(id string) (uint, error)

	// Returns amount of deleted deliveryProduct
	DeleteAll(filter DeliveryProductParams) (uint, error)

	// Not used, delivery-product should not be updated
	// Update(id string, newDeliveryProduct DeliveryProductParams) (DeliveryProductModel, error)
}
