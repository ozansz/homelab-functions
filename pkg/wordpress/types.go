package wordpress

import (
	"encoding/json"
	"fmt"
)

type Category struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
	Name string `json:"name"`
	Slug string `json:"slug"`

	Count       int    `json:"count,omitempty"`
	Description string `json:"description,omitempty"`
	Taxonomy    string `json:"taxonomy,omitempty"`
	Parent      int    `json:"parent,omitempty"`
	Links       struct {
		Self []struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
		Collection []struct {
			Href string `json:"href,omitempty"`
		} `json:"collection,omitempty"`
		About []struct {
			Href string `json:"href,omitempty"`
		} `json:"about,omitempty"`
	} `json:"_links,omitempty"`
}

type Comment struct {
	ID        int    `json:"id"`
	Post      int    `json:"post"`
	Parent    int    `json:"parent"`
	Author    int    `json:"author"`
	AuthorURL string `json:"author_url"`
	Date      string `json:"date"`
	Content   struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
	Link string `json:"link"`
}

type Page struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Link     string `json:"link"`
	Modified string `json:"modified"`
	Slug     string `json:"slug"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Title    struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Content struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
	Excerpt struct {
		Rendered string `json:"rendered"`
	} `json:"excerpt"`
	Author int `json:"author"`
	Parent int `json:"parent"`
}

type Post struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Modified string `json:"modified"`
	Slug     string `json:"slug"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Link     string `json:"link"`
	Title    struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Content struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
	Excerpt struct {
		Rendered string `json:"rendered"`
	} `json:"excerpt"`
	Author        int   `json:"author"`
	FeaturedMedia int   `json:"featured_media"`
	Categories    []int `json:"categories"`
	Tags          []int `json:"tags"`
}

type Tag struct {
	ID          int    `json:"id"`
	Count       int    `json:"count"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Taxonomy    string `json:"taxonomy"`
}

type User struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Description string            `json:"description"`
	Link        string            `json:"link"`
	Slug        string            `json:"slug"`
	AvatarURLs  map[string]string `json:"avatar_urls"`
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
