package product

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	var productParams ProductParams
	queryParams := r.URL.Query()

	if weight := queryParams.Get("weight_grams"); weight != "" {
		weightValue, err := strconv.ParseFloat(weight, 32)
		if err == nil {
			productParams.WeightGrams = float32(weightValue)
		} else {
			http.Error(w, "Invalid weight_grams value", http.StatusBadRequest)
			return
		}
	}

	if price := queryParams.Get("price"); price != "" {
		priceValue, err := strconv.ParseFloat(price, 32)
		if err == nil {
			productParams.Price = float32(priceValue)
		} else {
			http.Error(w, "Invalid price value", http.StatusBadRequest)
			return
		}
	}

	if name := queryParams.Get("name"); name != "" {
		productParams.Name = name
	}

	foundProducts, err := c.repository.GetAll(productParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	dtoFoundProjects := []ProductDTO{}
	for i := 0; i < len(foundProducts); i++ {
		dtoFoundProjects = append(dtoFoundProjects, foundProducts[i].ToDTO())
	}
	if err := json.NewEncoder(w).Encode(dtoFoundProjects); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *ProductController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c.repository.GetOne(id)
}

func (c *ProductController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	var productParams ProductParams
	queryParams := r.URL.Query()

	if weight := queryParams.Get("weight_grams"); weight != "" {
		weightValue, err := strconv.ParseFloat(weight, 32)
		if err == nil {
			productParams.WeightGrams = float32(weightValue)
		} else {
			http.Error(w, "Invalid weight_grams value", http.StatusBadRequest)
			return
		}
	}

	if price := queryParams.Get("price"); price != "" {
		priceValue, err := strconv.ParseFloat(price, 32)
		if err == nil {
			productParams.Price = float32(priceValue)
		} else {
			http.Error(w, "Invalid price value", http.StatusBadRequest)
			return
		}
	}

	if name := queryParams.Get("name"); name != "" {
		productParams.Name = name
	}

	deletedCount, err := c.repository.DeleteAll(productParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Since this won't be used elsewhere it can be created in here
	response := struct {
		DeletedCount uint `json:"deleted_count"`
	}{
		DeletedCount: deletedCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *ProductController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c.repository.DeleteOne(id)
}

func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c.repository.Update(id, ProductParams{})
}
