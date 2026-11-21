package vertex

type Vertex struct {
	Endpoint string `json:"endpoint"`
}

type EndpointUpdateRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
}
