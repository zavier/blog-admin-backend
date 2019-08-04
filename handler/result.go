package handler

type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func SuccessResult() Result {
	res := Result{
		Success: true,
		Data:    true,
		Message: "",
	}
	return res
}

func ErrorResult(msg string) Result {
	res := Result{
		Success: false,
		Data:    nil,
		Message: msg,
	}
	return res
}
