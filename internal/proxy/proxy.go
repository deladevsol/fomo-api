package proxy

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

//go:embed openapi.json
var specification []byte

var Routes = []string{
	"GET /v2/status",
	"GET /v2/users/handle/{handle}/theses",
	"GET /v2/users/handle/{handle}",
	"GET /v2/users/id/{id}",
	"GET /v2/users/wallet/{address}",
	"POST /v2/users/handle/{handle}/resolve",
	"GET /v2/users/handle/{handle}/pnl",
	"GET /v2/theses",
	"GET /v2/tokens/{tokenAddress}/theses",
	"GET /v2/users/id/{id}/theses",
	"GET /v2/users/id/{id}/tokens/{tokenAddress}/theses",
	"GET /v2/leaderboards/trader-pnl",
	"GET /v2/leaderboards/traders",
	"GET /v2/leaderboards/clans",
	"GET /v2/leaderboards/tokens/most-held",
	"GET /v2/leaderboards/tokens/trending",
	"GET /v2/leaderboards/tokens/graduated",
	"GET /v2/info",
	"GET /v2/stream",
	"GET /v2/users/search",
	"GET /v2/user/handle/{handle}",
	"GET /v2/user/wallet/{address}",
	"GET /v2/user/id/{id}",
	"POST /v2/user/handle/{handle}/resolve",
	"GET /v2/user/handle/{handle}/pnl",
	"GET /v2/thesis",
	"GET /v2/thesis/token/{tokenAddress}",
	"GET /v2/thesis/user/{id}",
	"GET /v2/thesis/user/{id}/token/{tokenAddress}",
	"GET /v2/leaderboard/traders-fomoscan",
	"GET /v2/leaderboard/traders",
	"GET /v2/leaderboard/clans",
	"GET /v2/leaderboard/tokens/most-held",
	"GET /v2/leaderboard/tokens/trending",
	"GET /v2/leaderboard/tokens/graduated",
	"GET /v2/me",
	"GET /v2/ws",
	"GET /v2/stats",
	"GET /v2/user/handle/{handle}/thesis",
	"GET /v2/thesis/handle/{handle}",
}

func New(base, key string) (http.Handler, error) {
	target, e := url.Parse(base)
	if e != nil || target.Host == "" || (target.Scheme != "https" && target.Scheme != "http") || target.User != nil || target.RawQuery != "" || target.Fragment != "" {
		return nil, fmt.Errorf("invalid upstream URL")
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 128
	transport.MaxIdleConnsPerHost = 32
	transport.MaxConnsPerHost = 64
	transport.ResponseHeaderTimeout = 30 * time.Second
	transport.IdleConnTimeout = 90 * time.Second
	p := &httputil.ReverseProxy{Transport: transport, FlushInterval: -1, Rewrite: func(r *httputil.ProxyRequest) {
		r.SetURL(target)
		r.Out.Host = target.Host
		r.Out.Header.Del("Authorization")
		r.Out.Header.Del("X-Api-Key")
		r.Out.Header.Del("Cookie")
		if key != "" {
			r.Out.Header.Set("Authorization", "Bearer "+key)
		}
	}, ErrorHandler: func(w http.ResponseWriter, r *http.Request, e error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(502)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "upstream unavailable"})
	}}
	mux := http.NewServeMux()
	for _, route := range Routes {
		mux.Handle(route, p)
	}
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"status\":\"ok\"}\n"))
	})
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(specification)
	})
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(docs))
	})
	return mux, nil
}

//go:embed docs.html
var docs string
