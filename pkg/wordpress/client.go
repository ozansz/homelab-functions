package wordpress

import (
	"context"
	"encoding/json"
	"fmt"
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
	// if content.Comments, err = c.GetComments(ctx); err != nil {
	// 	return nil, fmt.Errorf("failed to get comments: %w", err)
	// }
	// if content.Pages, err = c.GetPages(ctx); err != nil {
	// 	return nil, fmt.Errorf("failed to get pages: %w", err)
	// }
	// if content.Posts, err = c.GetPosts(ctx); err != nil {
	// 	return nil, fmt.Errorf("failed to get posts: %w", err)
	// }
	// if content.Tags, err = c.GetTags(ctx); err != nil {
	// 	return nil, fmt.Errorf("failed to get tags: %w", err)
	// }
	// if content.Users, err = c.GetUsers(ctx); err != nil {
	// 	return nil, fmt.Errorf("failed to get users: %w", err)
	// }
	return content, nil
}

func (c *Client) GetCategories(ctx context.Context) ([]Category, error) {
	categories := make([]Category, 0)

	for page := 1; ; page++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+categoriesPath, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		q := req.URL.Query()
		q.Add("per_page", fmt.Sprint(entitiesPerPage))
		q.Add("page", fmt.Sprint(page))
		req.URL.RawQuery = q.Encode()

		log.Printf("getting categories page %d, request: %s %s", page, req.Method, req.URL.String())

		res, err := c.cl.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}

		log.Printf("got response: %#v", res)
		defer res.Body.Close()

		switch res.StatusCode {
		case http.StatusOK:
			var pageCats []Category
			if err := json.NewDecoder(res.Body).Decode(&pageCats); err != nil {
				return nil, fmt.Errorf("failed to decode response body to []Category: %w", err)
			}
			categories = append(categories, pageCats...)
			if len(pageCats) < entitiesPerPage {
				log.Printf("got %d categories, less than %d, assuming that's all", len(pageCats), entitiesPerPage)
				return categories, nil
			}
		case http.StatusBadRequest:
			var errRes APIErrorResponse
			if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
				return nil, fmt.Errorf("failed to decode response body to APIErrorResponse: %w", err)
			}
			if errRes.IsInvalidPageNumber() {
				log.Printf("got invalid page number error, got total %d categories", len(categories))
				return categories, nil
			}
			return nil, fmt.Errorf("unexpected error: code: %s, message: %q", errRes.Code, errRes.Message)
		default:
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
	}
}
