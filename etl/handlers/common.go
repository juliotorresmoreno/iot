package handlers

import (
	"net/http"
)

type HttpBody struct {
	Message string
}

type HttpBodyData struct {
	Data  any
	Total int
}

type HttpErr struct {
	Status int
	Body   HttpBody
}

var ErrInternalServerError HttpErr = HttpErr{
	http.StatusInternalServerError,
	HttpBody{"ErrInternalServerError"},
}
var ErrBadRequest HttpErr = HttpErr{
	http.StatusBadRequest,
	HttpBody{"ErrBadRequest"},
}
