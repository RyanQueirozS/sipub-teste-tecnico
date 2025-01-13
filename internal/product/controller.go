package product

import (
	"encoding/json"
	"net/http"
)

type ProductController struct {
	// TODO
	// validatorService IValidatorService
	repository IProductRepository
}

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	var product ProductParams
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO Request validation should happen now
	createdProduct, err := c.repository.Create(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdProduct.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	c.repository.GetAll()
}

func (c *ProductController) GetOne(w http.ResponseWriter, r *http.Request) {
	c.repository.GetOne()
}

func (c *ProductController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	c.repository.DeleteAll()
}

func (c *ProductController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	c.repository.DeleteOne()
}

func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	c.repository.Update()
}
