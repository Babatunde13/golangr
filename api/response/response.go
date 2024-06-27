package response

type Response struct {
	Message 	string `json:"message" required:"true"`
	Data 		interface{} `json:"data" required:"false"`
	Error 		interface{} `json:"error" required:"false"`
	Status 		bool `json:"status" required:"true"`
}

func SuccessResponse(data interface{}, message string) Response {
	return Response{
		Message: message,
		Data: data,
		Status: true,
	}
}

func ErrorResponse(err error, message string) Response {
	return Response{
		Message: message,
		Error: err,
		Status: false,
	}
}
