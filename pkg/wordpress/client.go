package wordpress

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseAPIPath    = "/wp-json/wp/v2"
	categoriesPath = baseAPIPath + "/categories"
	commentsPath   = baseAPIPath + "/comments"
	pagesPath      = baseAPIPath + "/pages"
	postsPath      = baseAPIPath + "/posts"
	tagsPath       = baseAPIPath + "/tags"
	usersPath      = baseAPIPath + "/users"

	entitiesPerPage = 10
)

type Client struct {
	cl      *http.Client
	baseURL string
}

type NewClientOpt func(*Client)

func WithTimeout(t time.Duration) NewClientOpt {
	return func(c *Client) {
		c.cl.Timeout = t
	}
}

func NewClient(baseURL string, opts ...NewClientOpt) *Client {
	cl := &Client{
		cl: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
	for _, opt := range opts {
		opt(cl)
	}
	return cl
}

func (c *Client) GetAll(ctx context.Context) (*SiteContent, error) {
	var (
		err     error
		content *SiteContent = &SiteContent{}
	)
	if content.Categories, err = c.GetCategories(ctx); err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	if content.Comments, err = c.GetComments(ctx); err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	if content.Pages, err = c.GetPages(ctx); err != nil {
		return nil, fmt.Errorf("failed to get pages: %w", err)
	}
	if content.Posts, err = c.GetPosts(ctx); err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	if content.Tags, err = c.GetTags(ctx); err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	if content.Users, err = c.GetUsers(ctx); err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return content, nil
}

func (c *Client) GetCategories(ctx context.Context) ([]Category, error) {
	categories := make([]Category, 0)
	if err := c.paginatedRequest(ctx, c.baseURL+categoriesPath, func(b []byte) (int, error) {
		var cats []Category
		if err := json.Unmarshal(b, &cats); err != nil {
			return 0, fmt.Errorf("failed to unmarshal categories: %w", err)
		}
		categories = append(categories, cats...)
		return len(cats), nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

func (c *Client) GetComments(ctx context.Context) ([]Comment, error) {
	comments := make([]Comment, 0)
	if err := c.paginatedRequest(ctx, c.baseURL+commentsPath, func(b []byte) (int, error) {
		var comms []Comment
		if err := json.Unmarshal(b, &comms); err != nil {
			return 0, fmt.Errorf("failed to unmarshal comments: %w", err)
		}
		comments = append(comments, comms...)
		return len(comms), nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	return comments, nil
}

func (c *Client) GetPages(ctx context.Context) ([]Page, error) {
	pages := make([]Page, 0)
	if err := c.paginatedRequest(ctx, c.baseURL+pagesPath, func(b []byte) (int, error) {
		var pgs []Page
		if err := json.Unmarshal(b, &pgs); err != nil {
			return 0, fmt.Errorf("failed to unmarshal pages: %w", err)
		}
		pages = append(pages, pgs...)
		return len(pgs), nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get pages: %w", err)
	}
	return pages, nil
}

func (c *Client) GetPosts(ctx context.Context) ([]Post, error) {
	posts := make([]Post, 0)
	if err := c.paginatedRequest(ctx, c.baseURL+postsPath, func(b []byte) (int, error) {
		var pst []Post
		if err := json.Unmarshal(b, &pst); err != nil {
			return 0, fmt.Errorf("failed to unmarshal posts: %w", err)
		}
		posts = append(posts, pst...)
		return len(pst), nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	return posts, nil
}

func (c *Client) GetTags(ctx context.Context) ([]Tag, error) {
	tags := make([]Tag, 0)
	if err := c.paginatedRequest(ctx, c.baseURL+tagsPath, func(b []byte) (int, error) {
		var tgs []Tag
		if err := json.Unmarshal(b, &tgs); err != nil {
			return 0, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
		tags = append(tags, tgs...)
		return len(tgs), nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	return tags, nil
}

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	users := make([]User, 0)
	if err := c.paginatedRequest(ctx, c.baseURL+usersPath, func(b []byte) (int, error) {
		var usrs []User
		if err := json.Unmarshal(b, &usrs); err != nil {
			return 0, fmt.Errorf("failed to unmarshal users: %w", err)
		}
		users = append(users, usrs...)
		return len(usrs), nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

func (c *Client) paginatedRequest(ctx context.Context, path string, forEach func([]byte) (int, error)) error {
	for page := 1; ; page++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			return fmt.Errorf("failed to create HTTP request: %w", err)
		}

		q := req.URL.Query()
		q.Add("per_page", fmt.Sprint(entitiesPerPage))
		q.Add("page", fmt.Sprint(page))
		req.URL.RawQuery = q.Encode()

		log.Printf("paginated HTTP request page: %d, %s %s", page, req.Method, req.URL.String())

		res, err := c.cl.Do(req)
		if err != nil {
			return fmt.Errorf("failed to do HTTP request: %w", err)
		}

		log.Printf("got response: %#v", res)
		defer res.Body.Close()

		switch res.StatusCode {
		case http.StatusOK:
			var b []byte
			if b, err = io.ReadAll(res.Body); err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}
			entities, err := forEach(b)
			if err != nil {
				return fmt.Errorf("failed to process response body: %w", err)
			}
			if entities < entitiesPerPage {
				return nil
			}
		case http.StatusBadRequest:
			var errRes APIErrorResponse
			if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
				return fmt.Errorf("failed to decode error response: %w", err)
			}
			if errRes.IsInvalidPageNumber() {
				log.Printf("got invalid page number error")
				return nil
			}
			return fmt.Errorf("got bad request error: %s", errRes.Message)
		default:
			return fmt.Errorf("got unexpected status code: %d", res.StatusCode)
		}

		// REMOVE AFTER TESTING
		if page == 3 {
			return nil
		}
	}
}
