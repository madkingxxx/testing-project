package presenter

type HttpResponse struct {
	// HTTP status code
	StatusCode int `json:"status_code"`
	// HTTP status message
	StatusMessage string `json:"status_message"`
	// HTTP response body
	Body interface{} `json:"body"`
}

var (
	// HTTP 200 OK
	OK = &HttpResponse{
		StatusCode:    200,
		StatusMessage: "OK",
		Body:          nil,
	}
)
