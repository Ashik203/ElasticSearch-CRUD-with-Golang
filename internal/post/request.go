package post
type createRequest struct {
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}
type updateRequest struct {
	ID string
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}

type deleteRequest struct {
	ID string
}
type findRequest struct {
	ID string
}
