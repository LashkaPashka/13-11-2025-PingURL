package payload

type Request struct {
	LinksList []string `json:"links_list"`
}

type Response struct {
	Status string `json:"status"`
}