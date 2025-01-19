package payment

type IPaymentRepository interface {
	// Returns the created payment
	Create(PaymentParams) (PaymentModel, error)

	// Returns the found payments
	// NOTE: reusing the same type for a filter and a "constructor" is not
	// ideal at all, but it will save on code repetition
	GetAll(filter PaymentParams) ([]PaymentModel, error)

	// Returns the found payment
	GetOne(id string) (PaymentModel, error)

	// Returns amount of deleted payments
	DeleteOne(id string) (uint, error)

	// Returns amount of deleted payments
	DeleteAll(filter PaymentParams) (uint, error)

	// Cannot be updated after being created
	// Update(id string, newPayment PaymentParams) (PaymentModel, error)
}
