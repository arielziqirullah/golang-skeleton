package helper

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

type EmptyObj struct {
}

func BuildResponse(status bool, message string, code int, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: message,
		Code:    code,
		Data:    data,
	}

	return res
}

func BuildErrorResponse(message string, code int, data interface{}) Response {
	res := Response{
		Status:  false,
		Message: message,
		Code:    code,
		Data:    data,
	}

	return res
}
