package server

import (
	"encoding/json"
	"net/http"
)

type JsonResponse interface {
	GetHttpCode() int
	GetInternalCode() int
	GetJson() ([]byte, error)
	WriteHttpResponse(w http.ResponseWriter) (int, error)
}

func NewOkJsonResponse(data interface{}) JsonResponse {
	r := concreteJsonResponse{
		ResponseCode: http.StatusOK,

		Err:  nil,
		Data: data,
	}

	return r
}

func NewErrorJsonResponse(data interface{}, responseCode int, err error, errorCode int) JsonResponse {
	errStr := err.Error()
	r := concreteJsonResponse{
		ResponseCode: responseCode,
		Err:          &errStr,
		Data:         data,
		Code:         errorCode,
	}

	return r
}

type concreteJsonResponse struct {
	ResponseCode int         `json:"-"`
	Err          *string     `json:"error"`
	Code         int         `json:"code"`
	Data         interface{} `json:"data"`
}

func (r concreteJsonResponse) GetHttpCode() int {
	return r.ResponseCode
}

func (r concreteJsonResponse) GetInternalCode() int {
	return r.Code
}

func (r concreteJsonResponse) GetJson() ([]byte, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r concreteJsonResponse) WriteHttpResponse(w http.ResponseWriter) (int, error) {
	jsonData, err := r.GetJson()
	if err != nil {
		return 0, err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(r.GetHttpCode())
	return w.Write(jsonData)
}
