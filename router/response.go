package router

import (
	"encoding/json"
	"fortest/validator"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	Status  int                    `json:"status"`
	Success bool                   `json:"success"`
	Meta    map[string]interface{} `json:"meta"`
	Data    interface{}            `json:"data"`
}

func NewResponse(status int, meta map[string]interface{}, data interface{}) *Response {
	success := false
	if status >= 200 && status <= 299 {
		success = true
	}
	response := &Response{
		Status:  status,
		Success: success,
		Meta:    meta,
		Data:    data,
	}

	if response.Data == nil {
		response.Data = http.StatusText(status)
	}

	if v, ok := data.(error); ok {
		response.Data = v.Error()
	}
	return response
}

func Respond(w http.ResponseWriter, r *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Println(err)
	}
}

func PaginationMeta(total int64, offset int, limit int) map[string]interface{} {
	return map[string]interface{}{
		"total":  total,
		"offset": offset,
		"limit":  limit,
	}
}

func OK(w http.ResponseWriter, data interface{}) {
	Respond(w, NewResponse(http.StatusOK, nil, data))
}

func OKMeta(w http.ResponseWriter, meta map[string]interface{}, data interface{}) {
	Respond(w, NewResponse(http.StatusOK, meta, data))
}

func BadRequest(w http.ResponseWriter, data interface{}) {
	Respond(w, NewResponse(http.StatusBadRequest, nil, data))
}

func ValidationFailed(w http.ResponseWriter, err error) {
	data := strings.Builder{}
	for _, e := range validator.CheckValidationErrors(err) {
		data.Write([]byte(", " + e.Error()))
	}

	Respond(w, NewResponse(http.StatusUnprocessableEntity, nil, data.String()))
}

func BadRequestMeta(w http.ResponseWriter, meta map[string]interface{}, data interface{}) {
	Respond(w, NewResponse(http.StatusBadRequest, meta, data))
}

func Unauthorized(w http.ResponseWriter, data interface{}) {
	Respond(w, NewResponse(http.StatusUnauthorized, nil, data))
}

func Forbidden(w http.ResponseWriter, data interface{}) {
	Respond(w, NewResponse(http.StatusForbidden, nil, data))
}

func NotFound(w http.ResponseWriter, data interface{}) {
	Respond(w, NewResponse(http.StatusNotFound, nil, data))
}

func Internal(w http.ResponseWriter, data interface{}) {
	Respond(w, NewResponse(http.StatusInternalServerError, nil, data))
}
