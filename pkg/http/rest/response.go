package rest

type Meta struct {
	Total int `json:"total"`
}

type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type CollectionResponse struct {
	Meta  Meta        `json:"meta"`
	Links []Link      `json:"links"`
	Items interface{} `json:"items"`
}
