package model

type Links struct {
	Urls       []Url  `json:"urls"`
	Status     string `json:"status"`
	NumberLink string `json:"links_num"`
}

type Url struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	UrlName    string `json:"url_name"`
	Available  string `json:"available"`
}
