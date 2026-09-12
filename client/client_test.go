package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestThesisCursorAndNullableValues(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/thesis/user/user-1/token/TokenCaseSensitive" || r.URL.Query().Get("before") != "cursor+/=" {
			t.Errorf("incorrect URL: %s", r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":1,"hasMore":false,"nextBefore":null,"updatedAt":1720000000000,"items":[{"id":"comment-1","thesis":"hello","likeCount":null,"tokenAmount":0}]}`))
	}))
	defer s.Close()
	c, e := New(Config{BaseURL: s.URL})
	if e != nil {
		t.Fatal(e)
	}
	page, e := c.Theses(context.Background(), ThesisQuery{UserID: "user-1", TokenAddress: "TokenCaseSensitive", Before: "cursor+/="})
	if e != nil {
		t.Fatal(e)
	}
	if page.NextBefore != nil || page.Items[0].LikeCount != nil || page.Items[0].TokenAmount == nil || *page.Items[0].TokenAmount != 0 {
		t.Fatalf("null and zero were conflated: %+v", page)
	}
}

func TestErrorsPreserveStatusAndRetryAfter(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "12")
		w.WriteHeader(429)
		_, _ = w.Write([]byte(`{"error":"try later"}`))
	}))
	defer s.Close()
	c, _ := New(Config{BaseURL: s.URL})
	_, e := c.UserByHandle(context.Background(), "test")
	var apiErr *Error
	if !errors.As(e, &apiErr) || apiErr.StatusCode != 429 || apiErr.RetryAfter != "12" {
		t.Fatalf("error: %v", e)
	}
}

func TestRedirectDoesNotLeakKey(t *testing.T) {
	leaked := false
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { leaked = true }))
	defer target.Close()
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, target.URL, 302) }))
	defer source.Close()
	c, _ := New(Config{BaseURL: source.URL, APIKey: "private"})
	_, e := c.UserByHandle(context.Background(), "test")
	if e == nil || leaked {
		t.Fatal("redirect followed")
	}
}

func TestPathCannotChangeOrigin(t *testing.T) {
	c, _ := New(Config{})
	for _, path := range []string{"https://example.com/v2/thesis", "//example.com/v2/thesis", "/v2/pump/thesis"} {
		if e := c.Request(context.Background(), "GET", path, nil); e == nil {
			t.Fatalf("accepted %q", path)
		}
	}
}
