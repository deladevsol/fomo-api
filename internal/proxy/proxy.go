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
	"GET /v2/user/handle/{handle}", "GET /v2/user/id/{id}", "GET /v2/user/wallet/{address}", "POST /v2/user/handle/{handle}/resolve", "GET /v2/user/handle/{handle}/pnl",
	"GET /v2/thesis", "GET /v2/thesis/token/{tokenAddress}", "GET /v2/thesis/user/{id}", "GET /v2/thesis/user/{id}/token/{tokenAddress}",
	"GET /v2/leaderboard/traders-fomoscan", "GET /v2/leaderboard/traders", "GET /v2/leaderboard/clans", "GET /v2/leaderboard/tokens/most-held", "GET /v2/leaderboard/tokens/trending", "GET /v2/leaderboard/tokens/graduated", "GET /v2/me", "GET /v2/ws",
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

const docs = `<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Fomo API</title><style>body{font:16px system-ui;max-width:900px;margin:48px auto;padding:0 24px;background:#101318;color:#e4e8ee}a{color:#9dbbff}code{font:14px monospace}td{padding:12px;border-bottom:1px solid #303843}table{border-collapse:collapse;width:100%}</style><h1>Fomo API</h1><p>Free access to public Fomo data, powered by <a href="https://pootracker.app">PooTracker</a>.</p><p><a href="/openapi.json">OpenAPI schema</a> · <a href="https://github.com/deladevsol/fomo-api">Source</a></p><table id="routes"></table><script>fetch('/openapi.json').then(r=>r.json()).then(s=>{for(const[p,ops]of Object.entries(s.paths))for(const m of Object.keys(ops)){const tr=document.createElement('tr');for(const v of [m.toUpperCase(),p]){const td=document.createElement('td');td.textContent=v;tr.append(td)}document.getElementById('routes').append(tr)}})</script></html>`
