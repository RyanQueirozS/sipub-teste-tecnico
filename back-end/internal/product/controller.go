package product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ProductController struct {
	// TODO
	validator  ProductValidator
	repository IProductRepository
}

// Used for testing
func (c *ProductController) SetRepository(repo IProductRepository) {
	c.repository = repo
}

func NewProductController() *ProductController {
	return &ProductController{repository: NewMySQLproductRepository()}
}

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	var productParam ProductParams
	err := json.NewDecoder(r.Body).Decode(&productParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.validator.Validate(productParam); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdProduct, err := c.repository.Create(productParam)
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

	// First it checks to see if the values in the querystring are valid
	for key := range queryParams {
		if strings.ToLower(key) != "isactive" &&
			strings.ToLower(key) != "isdeleted" &&
			strings.ToLower(key) != "createdAt" &&
			strings.ToLower(key) != "weightgrams" &&
			strings.ToLower(key) != "price" &&
			strings.ToLower(key) != "name" {
			http.Error(w, fmt.Sprintf("Invalid query parameter: %s", key), http.StatusBadRequest)
			return
		}
	}

	// If the values are valid it will check what each value is
	if isActive := queryParams.Get("IsActive"); isActive != "" {
		isActiveValue, err := strconv.ParseBool(isActive)
		if err == nil {
			*productParams.IsActive = isActiveValue
		}
	}

	if isDeleted := queryParams.Get("IsDeleted"); isDeleted != "" {
		isDeletedValue, err := strconv.ParseBool(isDeleted)
		if err == nil {
			*productParams.IsDeleted = isDeletedValue
		}
	}

	if createdAt := queryParams.Get("CreatedAt"); createdAt != "" {
		*productParams.CreatedAt = createdAt
	}

	if weight := queryParams.Get("WeightGrams"); weight != "" {
		weightValue, err := strconv.ParseFloat(weight, 32)
		if err == nil {
			*productParams.WeightGrams = float32(weightValue)
		}
	}

	if price := queryParams.Get("Price"); price != "" {
		priceValue, err := strconv.ParseFloat(price, 32)
		if err == nil {
			*productParams.Price = float32(priceValue)
		}
	}

	if name := queryParams.Get("Name"); name != "" {
		*productParams.Name = name
	}

	// NOTE: Why don't I just parse it from the json? Well, if the json is nil,
	// then it will generate an error and if the json isn't but a field is,
	// that is another possible error. This way, although repetitive, will make
	// it simple to understand. Where as having multiple nested `if`s might not

	// It now passes the product param as a "filter" and gets the found products
	foundProducts, err := c.repository.GetAll(productParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dtoFoundProducts := []ProductDTO{}
	for i := 0; i < len(foundProducts); i++ {
		dtoFoundProducts = append(dtoFoundProducts, foundProducts[i].ToDTO())
	}
	// Returns the DTO products
	if err := json.NewEncoder(w).Encode(dtoFoundProducts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *ProductController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	product, err := c.repository.GetOne(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // The Id was not found but the request did go though
		return
	}
	if err := json.NewEncoder(w).Encode(product.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *ProductController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	var productParams ProductParams
	queryParams := r.URL.Query()

	if isActive := queryParams.Get("IsActive"); isActive != "" {
		isActiveValue, err := strconv.ParseBool(isActive)
		if err == nil {
			*productParams.IsActive = isActiveValue
		}
	}

	if isDeleted := queryParams.Get("IsDeleted"); isDeleted != "" {
		isDeletedValue, err := strconv.ParseBool(isDeleted)
		if err == nil {
			*productParams.IsDeleted = isDeletedValue
		}
	}

	if createdAt := queryParams.Get("CreatedAt"); createdAt != "" {
		*productParams.CreatedAt = createdAt
	}

	if weight := queryParams.Get("WeightGrams"); weight != "" {
		weightValue, err := strconv.ParseFloat(weight, 32)
		if err == nil {
			*productParams.WeightGrams = float32(weightValue)
		}
	}

	if price := queryParams.Get("Price"); price != "" {
		priceValue, err := strconv.ParseFloat(price, 32)
		if err == nil {
			*productParams.Price = float32(priceValue)
		}
	}

	if name := queryParams.Get("Name"); name != "" {
		if productParams.Name == nil {
			productParams.Name = new(string)
		}
		*productParams.Name = name
	}

	count, err := c.repository.DeleteAll(productParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(count); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *ProductController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	count, err := c.repository.DeleteOne(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Did go through, none found
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(count); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var productParams ProductParams
	err := json.NewDecoder(r.Body).Decode(&productParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product, err := c.repository.Update(id, productParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Did go through, none found
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
