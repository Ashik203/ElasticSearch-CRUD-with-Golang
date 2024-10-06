package post

type createRequest struct {
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}
type updateRequest struct {
	ID    string   
	Title string   `json:"title,omitempty"`
	Text  string   `json:"text,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

type deleteRequest struct {
	ID string
}
type findRequest struct {
	ID string
}
