package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	statusOk    = "Ok"
	statusError = "Error"
)

func Error(error string) Response {
	return Response{
		Status: statusError,
		Error:  error,
	}
}

func Ok() Response {
	return Response{
		Status: statusOk,
	}
}
