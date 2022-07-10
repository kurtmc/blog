package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DevToClient struct {
	url    string
	apiKey string
}

func NewClient(url, apiKey string) *DevToClient {
	return &DevToClient{
		url:    url,
		apiKey: apiKey,
	}
}

func (d *DevToClient) GetArticles() ([]Article, error) {
	req, err := http.NewRequest("GET", d.url+"/api/articles/me/all", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Key", os.Getenv("DEV_TO_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var target []Article

	err = json.NewDecoder(resp.Body).Decode(&target)

	return target, err
}

func (d *DevToClient) CreateNewArticle(articleRequest *ArticleRequest) (*Article, error) {
	body, err := json.Marshal(articleRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", d.url+"/api/articles", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Key", os.Getenv("DEV_TO_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var target Article

	err = json.NewDecoder(resp.Body).Decode(&target)

	return &target, err
}

func (d *DevToClient) UpdateArticle(id int, artcileRequest *ArticleRequest) (*Article, error) {
	body, err := json.Marshal(artcileRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal articleRequest: %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/articles/%d", d.url, id), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %v", err)
	}
	req.Header.Set("Api-Key", os.Getenv("DEV_TO_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute PUT request: %v", err)
	}
	defer resp.Body.Close()

	var target Article

	err = json.NewDecoder(resp.Body).Decode(&target)

	return &target, err
}

type Article struct {
	TypeOf                 string      `json:"type_of"`
	ID                     int         `json:"id"`
	Title                  string      `json:"title"`
	Description            string      `json:"description"`
	Published              bool        `json:"published"`
	PublishedAt            interface{} `json:"published_at"`
	Slug                   string      `json:"slug"`
	Path                   string      `json:"path"`
	URL                    string      `json:"url"`
	CommentsCount          int         `json:"comments_count"`
	PublicReactionsCount   int         `json:"public_reactions_count"`
	PageViewsCount         int         `json:"page_views_count"`
	PublishedTimestamp     string      `json:"published_timestamp"`
	BodyMarkdown           string      `json:"body_markdown"`
	PositiveReactionsCount int         `json:"positive_reactions_count"`
	CoverImage             interface{} `json:"cover_image"`
	TagList                []string    `json:"tag_list"`
	CanonicalURL           string      `json:"canonical_url"`
	ReadingTimeMinutes     int         `json:"reading_time_minutes"`
	User                   struct {
		Name            string      `json:"name"`
		Username        string      `json:"username"`
		TwitterUsername string      `json:"twitter_username"`
		GithubUsername  string      `json:"github_username"`
		WebsiteURL      interface{} `json:"website_url"`
		ProfileImage    string      `json:"profile_image"`
		ProfileImage90  string      `json:"profile_image_90"`
	} `json:"user"`
}

type ArticleRequest struct {
	Article struct {
		Title        string   `json:"title"`
		BodyMarkdown string   `json:"body_markdown"`
		Published    bool     `json:"published,omitempty"`
		Tags         []string `json:"tags,omitempty"`
	} `json:"article"`
}
