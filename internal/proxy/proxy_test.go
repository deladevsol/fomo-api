package proxy

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestRoutesAndRequestIsolation(t *testing.T) {
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Cookie") != "" || r.Header.Get("X-Api-Key") != "" || r.Header.Get("Authorization") != "Bearer operator-key" {
			t.Errorf("unexpected upstream headers: %v", r.Header)
		}
		w.Header().Set("X-Upstream", "yes")
		w.WriteHeader(202)
		_, _ = w.Write([]byte(`{"pending":true}`))
	}))
	defer up.Close()
	handler, e := New(up.URL, "operator-key")
	if e != nil {
		t.Fatal(e)
	}
	placeholder := regexp.MustCompile(`\{[^}]+\}`)
	for _, route := range Routes {
		parts := strings.SplitN(route, " ", 2)
		p := placeholder.ReplaceAllString(parts[1], "sample")
		r := httptest.NewRequest(parts[0], p+"?before=cursor", nil)
		r.Header.Set("Authorization", "Bearer caller-secret")
		r.Header.Set("Cookie", "session=private")
		r.Header.Set("X-Api-Key", "caller-secret")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != 202 || w.Header().Get("X-Upstream") != "yes" {
			t.Errorf("%s: %d %s", route, w.Code, w.Body.String())
		}
	}
	for _, path := range []string{"/v2/pump/thesis", "/admin/jobs", "/v2/unknown"} {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != 404 {
			t.Errorf("unexpected route %s: %d", path, w.Code)
		}
	}
}

func TestOpenAPIRoutesMatchProxy(t *testing.T) {
	var spec struct {
		Paths map[string]map[string]any `json:"paths"`
	}
	if e := json.Unmarshal(specification, &spec); e != nil {
		t.Fatal(e)
	}
	have := map[string]bool{}
	for _, r := range Routes {
		have[r] = true
	}
	for path, ops := range spec.Paths {
		if path == "/healthz" {
			continue
		}
		for method := range ops {
			if !have[strings.ToUpper(method)+" "+path] {
				t.Errorf("documented route missing: %s %s", method, path)
			}
		}
	}
}

func TestUpstreamFailureIsBoundedAndSanitized(t *testing.T) {
	handler, _ := New("http://127.0.0.1:1", "secret")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/v2/thesis", nil))
	body, _ := io.ReadAll(w.Result().Body)
	if w.Code != 502 || strings.Contains(string(body), "secret") {
		t.Fatalf("%d %s", w.Code, body)
	}
}
