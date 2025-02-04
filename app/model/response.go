package model

// Response is used for static shape json return
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

// BuildResponse method is to inject data value to dynamic success response
func BuildResponse(message string, data interface{}, err interface{}) Response {
	var errmess interface{} = nil
	if err != nil {
		errmess = err
	}
	res := Response{
		Status:  true,
		Message: message,
		Error:   errmess,
		Data:    data,
	}
	return res
}

// BuildErrorResponse method is to inject data value to dynamic failed response
func BuildErrorResponse(message string, err interface{}, data interface{}) Response {
	res := Response{
		Status:  false,
		Message: message,
		Error:   err,
		Data:    data,
	}
	return res
}
