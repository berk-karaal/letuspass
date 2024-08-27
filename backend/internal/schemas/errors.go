package schemas

type BadRequestResponse struct {
	Error string `json:"error" binding:"required"`
}

type NotFoundResponse struct {
	Error string `json:"error" binding:"required"`
}
