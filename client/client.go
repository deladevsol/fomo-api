// Package client provides access to the public PooTracker Fomo API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultURL = "https://fomo-public.pootracker.app"

type Client struct {
	base *url.URL
	http *http.Client
	key  string
}
type Config struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func New(cfg Config) (*Client, error) {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultURL
	}
	u, e := url.Parse(cfg.BaseURL)
	if e != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, fmt.Errorf("invalid API base URL")
	}
	u.Path = strings.TrimRight(u.Path, "/")
	h := cfg.HTTPClient
	if h == nil {
		h = &http.Client{Timeout: 30 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	}
	return &Client{base: u, http: h, key: cfg.APIKey}, nil
}

type Error struct {
	StatusCode int
	Message    string
	RetryAfter string
}

func (e *Error) Error() string { return fmt.Sprintf("API HTTP %d: %s", e.StatusCode, e.Message) }

func (c *Client) Request(ctx context.Context, method, path string, out any) error {
	if method != http.MethodGet && method != http.MethodPost {
		return fmt.Errorf("unsupported method")
	}
	u, e := url.Parse(path)
	if e != nil || u.IsAbs() || u.Host != "" || !strings.HasPrefix(path, "/v2/") || strings.HasPrefix(u.Path, "/v2/pump") {
		return fmt.Errorf("expected a non-pump /v2/ API path")
	}
	req, e := http.NewRequestWithContext(ctx, method, c.base.String()+path, nil)
	if e != nil {
		return e
	}
	req.Header.Set("Accept", "application/json")
	if c.key != "" {
		req.Header.Set("Authorization", "Bearer "+c.key)
	}
	r, e := c.http.Do(req)
	if e != nil {
		return e
	}
	defer r.Body.Close()
	body, e := io.ReadAll(io.LimitReader(r.Body, (16<<20)+1))
	if e != nil {
		return e
	}
	if len(body) > 16<<20 {
		return fmt.Errorf("API response exceeds 16 MiB")
	}
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		var v struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &v)
		if v.Error == "" {
			v.Error = http.StatusText(r.StatusCode)
		}
		return &Error{r.StatusCode, v.Error, r.Header.Get("Retry-After")}
	}
	if out == nil {
		return nil
	}
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	return d.Decode(out)
}

func (c *Client) UserByHandle(ctx context.Context, handle string) (User, error) {
	var u User
	e := c.Request(ctx, "GET", "/v2/users/handle/"+url.PathEscape(strings.TrimPrefix(handle, "@")), &u)
	return u, e
}
func (c *Client) UserByID(ctx context.Context, id string) (User, error) {
	var u User
	e := c.Request(ctx, "GET", "/v2/users/id/"+url.PathEscape(id), &u)
	return u, e
}
func (c *Client) UserByWallet(ctx context.Context, address string) (User, error) {
	var u User
	e := c.Request(ctx, "GET", "/v2/users/wallet/"+url.PathEscape(address), &u)
	return u, e
}
func (c *Client) Resolve(ctx context.Context, handle string) (User, error) {
	var u User
	e := c.Request(ctx, "POST", "/v2/users/handle/"+url.PathEscape(strings.TrimPrefix(handle, "@"))+"/resolve", &u)
	return u, e
}

type ThesisQuery struct{ UserID, TokenAddress, Before string }

func (c *Client) Theses(ctx context.Context, q ThesisQuery) (ThesisPage, error) {
	p := "/v2/theses"
	if q.UserID != "" {
		p = "/v2/users/id/" + url.PathEscape(q.UserID) + "/theses"
	}
	if q.TokenAddress != "" {
		if q.UserID != "" {
			p = "/v2/users/id/" + url.PathEscape(q.UserID) + "/tokens/" + url.PathEscape(q.TokenAddress) + "/theses"
		} else {
			p = "/v2/tokens/" + url.PathEscape(q.TokenAddress) + "/theses"
		}
	}
	if q.Before != "" {
		p += "?" + url.Values{"before": {q.Before}}.Encode()
	}
	var page ThesisPage
	e := c.Request(ctx, "GET", p, &page)
	return page, e
}
func (c *Client) Leaderboard(ctx context.Context, board string, q url.Values) (json.RawMessage, error) {
	switch board {
	case "traders", "trader-pnl", "traders-fomoscan", "clans", "tokens/most-held", "tokens/trending", "tokens/graduated":
	default:
		return nil, fmt.Errorf("unknown leaderboard")
	}
	if board == "traders-fomoscan" {
		board = "trader-pnl"
	}
	p := "/v2/leaderboards/" + board
	if len(q) > 0 {
		p += "?" + q.Encode()
	}
	var b json.RawMessage
	e := c.Request(ctx, "GET", p, &b)
	return b, e
}

// ThesesByHandle returns a page of archived theses without requiring a user ID.
func (c *Client) ThesesByHandle(ctx context.Context, handle, before string) (ThesisPage, error) {
	path := "/v2/users/handle/" + url.PathEscape(strings.TrimPrefix(strings.TrimSpace(handle), "@")) + "/theses"
	if before != "" {
		path += "?" + url.Values{"before": {before}}.Encode()
	}
	var page ThesisPage
	err := c.Request(ctx, "GET", path, &page)
	return page, err
}
