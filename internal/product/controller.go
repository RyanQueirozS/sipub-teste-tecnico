package product

import (
	"net/http"
)

type ProductController struct {
	// TODO
	// validatorService IValidatorService
	repository IProductRespository
}

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	// c.repository.Create()
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
