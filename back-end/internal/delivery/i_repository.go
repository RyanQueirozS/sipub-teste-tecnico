package delivery

type IDeliveryRepository interface {
	// Returns the created delivery
	Create(DeliveryParams) (DeliveryModel, error)

	// Returns the found deliveriees
	// NOTE: reusing the same type for a filter and a "constructor" is not
	// ideal at all, but it will save on code repetition
	GetAll(filter DeliveryParams) ([]DeliveryModel, error)

	// Returns the found delivery
	GetOne(id string) (DeliveryModel, error)

	// Returns amount of deleted deliveries
	DeleteOne(id string) (uint, error)

	// Returns amount of deleted deliveries
	DeleteAll(filter DeliveryParams) (uint, error)

	// Returns the updated delivery
	Update(id string, newDelivery DeliveryParams) (DeliveryModel, error)
}
