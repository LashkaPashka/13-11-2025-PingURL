package payload

type Request struct {
	LinksList []string `json:"links_list" validate:"required"`
}

type Response struct {
	Status string `json:"status"`
}