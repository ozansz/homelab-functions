package wordpress

import (
	"encoding/json"
	"fmt"
)

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

func (c *SiteContent) Marshal() (map[string][]byte, error) {
	comments, err := json.Marshal(c.Comments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal comments: %w", err)
	}
	pages, err := json.Marshal(c.Pages)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pages: %w", err)
	}
	posts, err := json.Marshal(c.Posts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal posts: %w", err)
	}
	categories, err := json.Marshal(c.Categories)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal categories: %w", err)
	}
	tags, err := json.Marshal(c.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}
	users, err := json.Marshal(c.Users)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal users: %w", err)
	}

	return map[string][]byte{
		"comments.json":   comments,
		"pages.json":      pages,
		"posts.json":      posts,
		"categories.json": categories,
		"tags.json":       tags,
		"users.json":      users,
	}, nil
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
