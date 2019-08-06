package common

import "net/http"

type Result struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

const (
	StatusUnauthorized        = http.StatusUnauthorized
	StatusOK                  = http.StatusOK
	StatusInternalServerError = http.StatusInternalServerError
	StatusBadRequest          = http.StatusBadRequest
)

func SuccessResult(data interface{}) Result {
	res := Result{
		Success: true,
		Code:    StatusOK,
		Data:    data,
		Message: "",
	}
	return res
}

func ErrorResult(code int, msg string) Result {
	res := Result{
		Success: false,
		Code:    code,
		Data:    nil,
		Message: msg,
	}
	return res
}
