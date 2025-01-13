// TODO
package product_test // According o go's https://pkg.go.dev/testing

import (
	"sipub-test/internal/product"
	"testing"

	"github.com/stretchr/testify/assert"
)

var inMemoryRepository = product.NewInMemoryProductRepository()

// I don't use a TestMain, but a setup function because each subtest (t.Run)
// will not be reset in a TestMain. TestMain only resets the test function
// itself (TestXxx). If we could define a new method for `testing.T` I could
// just add this to it to simplify workflow.
func beforeEachRepo() {
	// Each test should be independent of one another, so the inMemoryRepo will be reset
	inMemoryRepository = product.NewInMemoryProductRepository()
}

// Should not validade if a field is missing and should initialize field as 0 -
// if not specified - as validation is a responsibility of the controller.
// NOTE: I should have used random value generator for this, but because this
// is so simple, this should do it!
func TestRepositoryCreate(t *testing.T) {
	beforeEachRepo()
	t.Run("CreateProductWithOnlyWeight", func(t *testing.T) {
		product := inMemoryRepository.Create(product.ProductParams{
			WeightGrams: 100.0,
		})
		// We use EqualValues because golang is really typesafe and requires the
		// type to be the same when just comparing equality. The function returns
		// float32 and 0.0, for example, is float64, so it can't be equal
		// although the value is the same
		assert.EqualValues(t, product.GetWeight(), 100, "Should initialize field as zero if not specified")
		assert.EqualValues(t, product.GetPrice(), 0, "Should initialize field as zero if not specified")
	})

	beforeEachRepo()
	t.Run("CreateProductWithOnlyPrice", func(t *testing.T) {
		product := inMemoryRepository.Create(product.ProductParams{
			Price: 100.0,
		})
		assert.EqualValues(t, product.GetWeight(), 0, "Should initialize field as zero if not specified")
		assert.EqualValues(t, product.GetPrice(), 100, "Should initialize field as zero if not specified")
	})

	beforeEachRepo()
	t.Run("CreateProductWithWeightAndPrice", func(t *testing.T) {
		product := inMemoryRepository.Create(product.ProductParams{
			WeightGrams: 100.0,
			Price:       25.5,
		})
		assert.EqualValues(t, product.GetPrice(), 25.5, "Price should be set correctly")
		assert.EqualValues(t, product.GetWeight(), 100, "Weight should be set correctly")
	})

	beforeEachRepo()
	t.Run("CreateProductWithDefaults", func(t *testing.T) {
		product := inMemoryRepository.Create(product.ProductParams{})
		assert.EqualValues(t, product.GetPrice(), 0, "Should initialize price as zero by default")
		assert.EqualValues(t, product.GetWeight(), 0, "Should initialize weight as zero by default")
	})
}
