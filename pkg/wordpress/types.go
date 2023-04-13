package wordpress

type Comment struct {
	// TODO
}

type Page struct {
	// TODO
}

type Post struct {
	// TODO
}

type Category struct {
	// TODO
}

type Tag struct {
	// TODO
}

type User struct {
	// TODO
}

type SiteContent struct {
	Comments   []Comment
	Pages      []Page
	Posts      []Post
	Categories []Category
	Tags       []Tag
	Users      []User
}

type APIErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Status int `json:"status"`
	} `json:"data,omitempty"`
}

func (e *APIErrorResponse) IsInvalidPageNumber() bool {
	return e.Code == "rest_post_invalid_page_number"
}
