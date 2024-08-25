package schemas

type BadRequestResponse struct {
	Error string `json:"error"`
}

type NotFoundResponse struct {
	Error string `json:"error"`
}
