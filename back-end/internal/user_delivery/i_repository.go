package user_delivery

type IUserDeliveryRepository interface {
	// Returns the created userDelivery
	Create(UserDeliveryParams) (UserDeliveryModel, error)

	// Returns the found userDeliverys
	// NOTE: reusing the same type for a filter and a "constructor" is not
	// ideal at all, but it will save on code repetition
	GetAll(filter UserDeliveryParams) ([]UserDeliveryModel, error)

	// Returns the found userDelivery
	GetOne(id string) (UserDeliveryModel, error)

	// Returns amount of deleted userDelivery
	DeleteOne(id string) (uint, error)

	// Returns amount of deleted userDelivery
	DeleteAll(filter UserDeliveryParams) (uint, error)

	// Not used, delivery-product should not be updated
	// Update(id string, newUserDelivery UserDeliveryParams) (UserDeliveryModel, error)
}
