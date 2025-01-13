package internal

import "net/http"

type IController interface {
	Create(http.ResponseWriter, *http.Request)

	GetAll(http.ResponseWriter, *http.Request)

	GetOne(http.ResponseWriter, *http.Request)

	DeleteOne(http.ResponseWriter, *http.Request)

	DeleteAll(http.ResponseWriter, *http.Request)

	Update(http.ResponseWriter, *http.Request)
}
