package bind

type Message struct {
	Message string `json:"message"`
}

type DataList struct {
	Total int64 `json:"total"`
	Data  any   `json:"data"`
}
