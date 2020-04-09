package rest

type Meta struct {
	Total int
}

type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type CollectionResponse struct {
	Meta  Meta        `json:"Meta"`
	Links []Link      `json:"links"`
	Items interface{} `json:"items"`
}
