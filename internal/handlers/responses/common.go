package responses

type EmptyResponse struct{}

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
